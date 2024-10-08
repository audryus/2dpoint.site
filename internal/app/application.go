package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/audryus/2dpoint.site/internal/config"
	"github.com/audryus/2dpoint.site/internal/controller"
	"github.com/audryus/2dpoint.site/internal/domain/memo"
	"github.com/audryus/2dpoint.site/internal/domain/memo/text"
	"github.com/audryus/2dpoint.site/internal/domain/memo/url"
	"github.com/audryus/2dpoint.site/internal/server"
	"github.com/audryus/2dpoint.site/internal/usecase"
	"github.com/audryus/2dpoint.site/pkg/database/cockroach"
	"github.com/audryus/2dpoint.site/pkg/database/etcd"
	"github.com/audryus/2dpoint.site/pkg/logger"
)

func Run() {
	logger := logger.New()

	cfg, err := config.New(logger)

	if err != nil {
		log.Fatal(err)
	}

	etcdClient, err := etcd.New(cfg, logger)
	db, err := cockroach.New(cfg, logger)

	urlRepo := url.NewUrlRepo(etcdClient)
	urlCreateS := url.NewCreateUrlService(urlRepo)
	urlGetService := url.NewGetUrlService(urlRepo)

	textRepo := text.NewTextRepo(etcdClient, db)
	createTextService := text.NewCreateTextService(textRepo)
	getTextService := text.NewGetTextService(textRepo)

	memoRepo := memo.NewMemoRepo(etcdClient)

	usecases := usecase.NewUseCases(usecase.Deps{
		CreateMemoService: memo.NewCreateMemoService(memoRepo, createTextService, urlCreateS),
		FetchMemoService:  memo.NewFetchService(memoRepo, getTextService, urlGetService),
	})

	controller := controller.NewController(usecases)
	srv := server.NewServer(cfg, controller.Init(cfg))

	go func() {
		if err := srv.Run(); err != nil {
			logger.Error(err, "error occurred while running http server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	const timeout = 5 * time.Second
	logger.Info("shutting down server...")
	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := srv.Stop(ctx); err != nil {
		logger.Error(err, "fiber app error")
	}
	logger.Info("server stopped")
	if err := etcdClient.Close(); err != nil {
		logger.Error(err, "failed to etcd client")
	}
	logger.Info("etcd disconnected")

	db.Close()
	logger.Info("cockroach disconnected")
}
