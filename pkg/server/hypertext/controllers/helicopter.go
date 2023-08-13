package controllers

import (
	"code.smartsheep.studio/atom/bedrock/pkg/server/datasource/models"
	"code.smartsheep.studio/atom/bedrock/pkg/server/hypertext/hyperutils"
	"code.smartsheep.studio/atom/bedrock/pkg/server/services"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"github.com/spf13/viper"
	"time"
)

type HeLiCoPtErController struct {
	cop *services.HeLiCoPtErService
}

func NewHeLiCoPtErController(cop *services.HeLiCoPtErService) *HeLiCoPtErController {
	return &HeLiCoPtErController{cop}
}

func (v *HeLiCoPtErController) Map(router *fiber.App) {
	router.Post("/cgi/subapps/:name", v.subappCommit)
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
