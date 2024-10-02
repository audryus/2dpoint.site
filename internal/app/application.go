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
	domain "github.com/audryus/2dpoint.site/internal/domain/memo"
	"github.com/audryus/2dpoint.site/internal/server"
	"github.com/audryus/2dpoint.site/internal/usecase"
	"github.com/audryus/2dpoint.site/pkg/database/etcd"
)

func Run() {
	cfg, err := config.New()

	if err != nil {
		log.Fatal(err)
	}
	etcdClient, err := etcd.NewClient(cfg)

	createMemoRepoEtcd := domain.NewCreateMemoRepoEtcd(etcdClient)
	getMemoRepoEtcd := domain.NewGetMemoRepoEtcd(etcdClient)

	usecases := usecase.NewUseCases(usecase.Deps{
		CreateMemoService: domain.NewCreateMemoService(createMemoRepoEtcd, getMemoRepoEtcd),
		GetMemoService:    domain.NewGetMemoService(getMemoRepoEtcd, createMemoRepoEtcd),
	})

	controller := controller.NewController(usecases)
	srv := server.NewServer(cfg, controller.Init(cfg))

	go func() {
		if err := srv.Run(); err != nil {
			log.Fatal("error occurred while running http server: %s\n", err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	const timeout = 5 * time.Second
	log.Println("shutting down server...")
	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := srv.Stop(ctx); err != nil {
		log.Fatal("failed to stop server: %v", err)
	}
	etcdClient.Close()

}
