package controller

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/audryus/2dpoint.site/internal/config"
	"github.com/audryus/2dpoint.site/internal/usecase"
	"github.com/audryus/2dpoint.site/pkg/logger"
	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/template/django/v3"
)

type Controller struct {
	usecases usecase.UseCases
}

func NewController(usecases usecase.UseCases) Controller {
	return Controller{
		usecases,
	}
}

func (ct Controller) Init(cfg config.Config) *fiber.App {
	engine := django.NewFileSystem(http.Dir("./internal/views"), ".html")

	app := fiber.New(fiber.Config{
		ServerHeader: cfg.Server.Header,
		AppName:      cfg.App.Name + " " + cfg.App.Version,
		Views:        engine,
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
	})

	fiberLogger := logger.New()
	app.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: fiberLogger.Core(),
	}))

	app.Get("/:memo", func(c *fiber.Ctx) error {
		uc := ct.usecases.GetMemo
		m := c.Params("memo")
		memo, err := uc.Get(m)
		if err != nil {
			return c.Status(404).SendString(err.Error())
		}
		conentType := c.Get("Content-Type", "text/html")

		switch conentType {
		case "application/json":
			return c.JSON(memo)
		}

		switch memo.Kind {
		case "text":
			return c.SendString(memo.Content)
		case "url":
			return c.Redirect(memo.Content)
		}

		return c.Status(404).SendString("something is wrong, i can feel ...")
	})

	app.Use(compress.New())

	static := fiber.Static{
		Compress:      true,
		MaxAge:        86400,
		CacheDuration: time.Hour,
	}

	app.Static("/public/favicon", "./internal/views/public/favicon", static)
	app.Static("/public/css", "./internal/views/public/css", static)
	app.Static("/public/fonts", "./internal/views/public/fonts", static)
	app.Static("/public/images", "./internal/views/public/images", static)
	app.Static("/public/js", "./internal/views/public/js", static)
	app.Static("/public/webfonts", "./internal/views/public/webfonts", static)

	app.Use(cache.New())

	app.Get("/", func(c *fiber.Ctx) error {
		//c.Set("Cache-Control", "private, max-age=86400")

		return c.Render("index", fiber.Map{})
	})

	app.Post("/", func(c *fiber.Ctx) error {
		uc := ct.usecases.CreateMemo

		req := new(Request)

		if err := c.BodyParser(req); err != nil {
			return c.Status(400).SendString("lero " + err.Error())
		}

		createdMemo, err := uc.Create(req.Content)

		if err != nil {
			return c.Status(400).SendString(err.Error())
		}
		createdMemo.Text.ID = cfg.Server.Addr + "/" + createdMemo.Text.ID
		for i := range createdMemo.Urls {
			createdMemo.Urls[i].ID = cfg.Server.Addr + "/" + createdMemo.Urls[i].ID
		}

		conentType := strings.Split(c.Get("Content-Type", "text/html"), ";")[0]

		switch conentType {
		case "application/json":
			return c.JSON(createdMemo)
		}

		return c.Render("memo", fiber.Map{
			"text": createdMemo.Text,
			"urls": createdMemo.Urls,
		})
	})

	return app
}

type Request struct {
	Content string `json:"content" form:"content"`
}
