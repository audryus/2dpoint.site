package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/audryus/2dpoint.site/internal/config"
	"github.com/audryus/2dpoint.site/internal/server"
)

func Run() {
	cfg, err := config.New()

	if err != nil {
		log.Fatal(err)
	}

	srv := server.NewServer(cfg)
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

	/*if err := mongoClient.Disconnect(context.Background()); err != nil {
		logger.Error(err.Error())
	}*/

}
