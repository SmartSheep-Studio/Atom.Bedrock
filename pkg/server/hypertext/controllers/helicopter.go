package controllers

import (
	"code.smartsheep.studio/atom/bedrock/pkg/server/hypertext/hyperutils"
	"code.smartsheep.studio/atom/bedrock/pkg/server/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"github.com/samber/lo"
	"time"
)

type HeLiCoPtErController struct {
	cop *services.HeLiCoPtErService
}

func NewHeLiCoPtErController(cop *services.HeLiCoPtErService) *HeLiCoPtErController {
	return &HeLiCoPtErController{cop}
}

func (v *HeLiCoPtErController) Map(router *fiber.App) {
	router.Get("/cgi/subapps/:name", v.subappRewrite)
	router.Post("/cgi/subapps", v.subappCommit)
}

func (v *HeLiCoPtErController) subappRewrite(c *fiber.Ctx) error {
	if app, ok := lo.Find(v.cop.Apps, func(item *services.HeLiCoPtErSubApp) bool {
		return item.Manifest.Name == c.Params("name")
	}); ok {
		if app.ExposedOptions == nil {
			return fiber.NewError(fiber.StatusBadGateway, "the app isn't exposed")
		} else {
			return proxy.Forward(app.ExposedOptions.URL)(c)
		}
	} else {
		return fiber.NewError(fiber.StatusBadGateway, "couldn't find app with your provided name")
	}
}

func (v *HeLiCoPtErController) subappCommit(c *fiber.Ctx) error {
	var req struct {
		URL string `json:"url"`
	}

	if err := hyperutils.BodyParser(c, &req); err != nil {
		return err
	}

	if app, ok := lo.Find(v.cop.Apps, func(item *services.HeLiCoPtErSubApp) bool {
		return item.Manifest.Name == c.Params("name")
	}); ok {
		if app.ExposedOptions == nil {
			app.ExposedOptions = &services.HeLiCoPtErExposedOptions{
				URL: req.URL,
			}
		}

		app.LastHealthyAt = lo.ToPtr(time.Now())

		return c.JSON(app)
	} else {
		return fiber.NewError(fiber.StatusBadRequest, "couldn't find app with your provided name")
	}
}
