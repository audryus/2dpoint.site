package controller

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/audryus/2dpoint.site/internal/config"
	"github.com/audryus/2dpoint.site/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/gofiber/fiber/v2/utils"
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
	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.Server.Addr,
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
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
		return c.Redirect(memo.Url)
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

	app.Use(cache.New())
	app.Use(encryptcookie.New(encryptcookie.Config{
		Key:    "LtceZ5qQJffGAJkzNTD1OE8Uq1WOhi4OmIWz+ciQyDg=",
		Except: []string{"2dpoint_token"},
	}))

	app.Use(csrf.New(csrf.Config{
		KeyLookup:      "header:X-Csrf-Token",
		CookieName:     "2dpoint_token",
		CookieSameSite: "Lax",
		Expiration:     1 * time.Hour,
		KeyGenerator:   utils.UUIDv4,
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		//c.Set("Cache-Control", "private, max-age=86400")

		csrf := c.Cookies("2dpoint_token", "")
		return c.Render("index", fiber.Map{
			"2dpoint_token": csrf,
		})
	})

	app.Post("/", func(c *fiber.Ctx) error {
		uc := ct.usecases.CreateMemo

		req := new(Request)

		if err := c.BodyParser(req); err != nil {
			return c.Status(400).SendString("lero " + err.Error())
		}

		createdMemo, err := uc.Create(req.URL, req.Type)

		if err != nil {
			return c.Status(400).SendString("lera " + err.Error())
		}

		conentType := strings.Split(c.Get("Content-Type", "text/html"), ";")[0]

		switch conentType {
		case "application/json":
			return c.JSON(createdMemo)
		}

		csrf := c.Cookies("2dpoint_token", "")

		return c.Render("index", fiber.Map{
			"memo":          createdMemo,
			"minimezed":     cfg.Server.Addr + "/" + createdMemo.ID,
			"2dpoint_token": csrf,
		})
	})

	return app
}

type Request struct {
	Type string `json:"type" form:"type"`
	URL  string `json:"url" form:"url"`
}
