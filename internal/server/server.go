package server

import (
	"context"
	"net/http"

	"github.com/goccy/go-json"

	"github.com/audryus/2dpoint.site/internal/config"
	"github.com/audryus/2dpoint.site/internal/controller"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/django/v3"
)

type Controller interface {
	Init()
}

type HttpServer struct {
	fiber      *fiber.App
	cfg        config.Config
	controller Controller
}

func NewServer(cfg config.Config) *HttpServer {
	engine := django.NewFileSystem(http.Dir("./internal/views"), ".html")
	fiber := fiber.New(fiber.Config{
		ServerHeader: "Fiber",
		AppName:      "Test App v1.0.1",
		Views:        engine,
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
	})

	fiber.Static("/public/favicon", "./internal/views/public/favicon")
	fiber.Static("/public/css", "./internal/views/public/css")
	fiber.Static("/public/fonts", "./internal/views/public/fonts")
	fiber.Static("/public/images", "./internal/views/public/images")
	fiber.Static("/public/js", "./internal/views/public/js")

	return &HttpServer{
		fiber:      fiber,
		controller: controller.NewController(cfg, fiber),
		cfg:        cfg,
	}
}

func (s *HttpServer) Run() error {

	s.controller.Init()
	return s.fiber.Listen(s.cfg.HTTP.Addr)
}

func (s *HttpServer) Stop(ctx context.Context) error {
	return s.fiber.Shutdown()
}
