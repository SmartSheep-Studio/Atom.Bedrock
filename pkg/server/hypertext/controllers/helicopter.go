package controllers

import (
	"code.smartsheep.studio/atom/bedrock/pkg/server/datasource/models"
	"code.smartsheep.studio/atom/bedrock/pkg/server/hypertext/hyperutils"
	"code.smartsheep.studio/atom/bedrock/pkg/server/services"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"github.com/samber/lo"
	"github.com/spf13/viper"
	"strings"
	"time"
)

type HeLiCoPtErController struct {
	cop *services.HeLiCoPtErService
}

func NewHeLiCoPtErController(cop *services.HeLiCoPtErService) *HeLiCoPtErController {
	return &HeLiCoPtErController{cop}
}

func (v *HeLiCoPtErController) Map(router *fiber.App) {
	router.All("/srv/subapps/:name/*", v.subappRewrite)
	router.Post("/cgi/subapps/:name", v.subappCommit)
}

func (v *HeLiCoPtErController) subappRewrite(c *fiber.Ctx) error {
	if app, ok := lo.Find(v.cop.Apps, func(item *models.SubApp) bool {
		return item.Manifest.Name == c.Params("name")
	}); ok {
		if app.ExposedOptions == nil {
			return fiber.NewError(fiber.StatusBadGateway, "the app isn't exposed")
		} else {
			prefix := fmt.Sprintf("/srv/subapps/%s", c.Params("name"))
			url := strings.ReplaceAll(string(c.Request().URI().Path()), prefix, "") + "?" + string(c.Request().URI().QueryString())

			return proxy.Forward(app.ExposedOptions.URL + url)(c)
		}
	} else {
		return fiber.NewError(fiber.StatusBadGateway, "couldn't find app with your provided name")
	}
}

func (v *HeLiCoPtErController) subappCommit(c *fiber.Ctx) error {
	var req struct {
		URL   string                     `json:"url"`
		Pages []models.SubAppExposedPage `json:"pages"`
	}

	if err := hyperutils.BodyParser(c, &req); err != nil {
		return err
	}

	if app, ok := lo.Find(v.cop.Apps, func(item *models.SubApp) bool {
		return item.Manifest.Name == c.Params("name")
	}); ok {
		if app.ExposedOptions == nil {
			app.ExposedOptions = &models.SubAppExposedOptions{
				URL:   req.URL,
				Pages: req.Pages,
			}
		}

		app.LastHealthyAt = lo.ToPtr(time.Now())

		return c.JSON(fiber.Map{
			"others":        v.cop.Apps,
			"configuration": viper.AllSettings(),
		})
	} else {
		return fiber.NewError(fiber.StatusBadRequest, "couldn't find app with your provided name")
	}
}
