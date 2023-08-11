package controllers

import (
	"code.smartsheep.studio/atom/bedrock/pkg/server/hypertext/middlewares"
	"code.smartsheep.studio/atom/bedrock/pkg/server/services"
	"github.com/gofiber/fiber/v2"
	"strings"

	"github.com/spf13/viper"
)

type StatusController struct {
	metrics    *services.MetricsService
	gatekeeper *middlewares.AuthMiddleware
}

func NewStatusController(metrics *services.MetricsService, gatekeeper *middlewares.AuthMiddleware) *StatusController {
	ctrl := &StatusController{metrics, gatekeeper}
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
		"nav": []any{},
	})
}
