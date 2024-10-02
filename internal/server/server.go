package server

import (
	"context"

	"github.com/audryus/2dpoint.site/internal/config"
	"github.com/gofiber/fiber/v2"
)

type Controller interface {
	Init()
}

type HttpServer struct {
	fiber *fiber.App
	cfg   config.Config
}

func NewServer(cfg config.Config, fiber *fiber.App) *HttpServer {
	return &HttpServer{
		fiber: fiber,
		cfg:   cfg,
	}
}

func (s *HttpServer) Run() error {
	return s.fiber.Listen(s.cfg.HTTP.Addr)
}

func (s *HttpServer) Stop(ctx context.Context) error {
	return s.fiber.Shutdown()
}
