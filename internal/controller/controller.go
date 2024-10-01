package controller

import (
	"fmt"

	"github.com/audryus/2dpoint.site/internal/config"
	"github.com/audryus/2dpoint.site/internal/usecase"
	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	cfg   config.Config
	fiber *fiber.App
}

func NewController(cfg config.Config, fiber *fiber.App) Controller {
	return Controller{
		cfg:   cfg,
		fiber: fiber,
	}
}

func (c Controller) Init() {
	app := c.fiber

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{})
	})

	app.Get("/:memo", func(c *fiber.Ctx) error {
		m := c.Params("memo")
		memo, err := usecase.GetMemo(m)

		if err != nil {
			return c.Status(404).SendString(err.Error())
		}
		if memo.Type == "URL" {
			return c.Redirect(memo.Url)
		}

		return c.Redirect(memo.Url)
	})

	app.Post("/", func(c *fiber.Ctx) error {
		req := new(Request)

		if err := c.BodyParser(req); err != nil {
			return c.Status(400).SendString(err.Error())
		}

		createdMemo, err := usecase.CreateMemo(req.URL, req.Type)

		if err != nil {
			return c.Status(400).SendString(err.Error())
		}

		fmt.Printf("%+v\n", createdMemo)

		return c.Render("index", fiber.Map{})
	})
}

type Request struct {
	Type string `json:"type" form:"type"`
	URL  string `json:"url" form:"url"`
}
