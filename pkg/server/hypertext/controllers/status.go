package controllers

import (
	"code.smartsheep.studio/atom/bedrock/pkg/server/datasource/models"
	"code.smartsheep.studio/atom/bedrock/pkg/server/hypertext/hyperutils"
	"code.smartsheep.studio/atom/bedrock/pkg/server/hypertext/middlewares"
	"code.smartsheep.studio/atom/bedrock/pkg/server/services"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"strings"

	"github.com/spf13/viper"
)

type StatusController struct {
	metrics    *services.MetricsService
	cop        *services.HeLiCoPtErService
	gatekeeper *middlewares.AuthMiddleware
}

func NewStatusController(metrics *services.MetricsService, gatekeeper *middlewares.AuthMiddleware, cop *services.HeLiCoPtErService) *StatusController {
	ctrl := &StatusController{metrics, cop, gatekeeper}
	return ctrl
}

func (ctrl *StatusController) Map(router *fiber.App) {
	router.Get("/api", ctrl.status)
	router.Get("/api/info", ctrl.information)
}

func (ctrl *StatusController) status(c *fiber.Ctx) error {
	if ctrl.metrics.IsReady {
		return c.JSON(fiber.Map{
			"status": "ready",
			"startup": fiber.Map{
				"at":    ctrl.metrics.StartAt,
				"usage": ctrl.metrics.StartCost,
			},
		})
	} else {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"status": "starting",
			"startup": fiber.Map{
				"at":    ctrl.metrics.StartAt,
				"usage": ctrl.metrics.StartCost,
			},
		})
	}
}

func (ctrl *StatusController) information(c *fiber.Ctx) error {
	if !ctrl.metrics.IsReady {
		return fiber.NewError(fiber.StatusServiceUnavailable, "hypertext isn't prepared yet")
	}

	var nav []map[string]any
	for _, app := range ctrl.cop.Apps {
		nav = append(
			nav,
			lo.Map(app.ExposedOptions.Pages, func(item models.SubAppExposedPage, index int) map[string]any {
				v := hyperutils.CovertStructToMap(item)
				v["to"] = fmt.Sprintf(
					"%s%s",
					app.ExposedOptions.URL,
					v["path"],
				)

				return v
			})...,
		)
	}

	firmware := strings.Split(c.App().Config().AppName, " ")
	return c.JSON(fiber.Map{
		"debug":            viper.GetBool("general.debug"),
		"firmware":         strings.Join(firmware[:len(firmware)-1], " "),
		"firmware_version": firmware[len(firmware)-1],
		"manifest": fiber.Map{
			"name":        viper.GetString("general.name"),
			"description": viper.GetString("general.description"),
		},
		"http_server": fiber.Map{
			"network":       c.App().Config().Network,
			"prefork":       c.App().Config().Prefork,
			"max_body_size": c.App().Config().BodyLimit,
		},
		"nav": nav,
	})
}
