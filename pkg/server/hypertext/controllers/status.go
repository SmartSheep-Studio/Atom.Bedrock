package controllers

import (
	"code.smartsheep.studio/atom/bedrock/pkg/server/datasource/models"
	"code.smartsheep.studio/atom/bedrock/pkg/server/hypertext/hyperutils"
	"code.smartsheep.studio/atom/bedrock/pkg/server/hypertext/middlewares"
	"code.smartsheep.studio/atom/bedrock/pkg/server/services"
	"encoding/json"
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

	var pages []map[string]any
	for _, app := range ctrl.cop.Apps {
		if app.ExposedOptions != nil {
			pages = append(
				pages,
				lo.Map(app.ExposedOptions.Pages, func(item models.SubAppExposedPage, index int) map[string]any {
					v := hyperutils.CovertStructToMap(item)
					v["to"] = fmt.Sprintf(
						"%s/srv/subapps/%s%s",
						viper.GetString("base_url"),
						app.Manifest.Name,
						v["path"],
					)

					return v
				})...,
			)
		}
	}

	var nav []map[string]any
	var navIn []any
	raw, _ := json.Marshal(viper.Get("helicopter.nav.items"))
	_ = json.Unmarshal(raw, &navIn)

	for _, item := range lo.Map(navIn, func(item any, index int) map[string]any {
		return hyperutils.CovertStructToMap(item)
	}) {
		var build func(item map[string]any) map[string]any

		build = func(item map[string]any) map[string]any {
			if item["to"] != nil {
				return item
			}

			app, ok := lo.Find(ctrl.cop.Apps, func(v *models.SubApp) bool {
				if v.ExposedOptions == nil {
					return false
				}
				for _, page := range v.ExposedOptions.Pages {
					if page.Name == item["name"] {
						return true
					}
				}
				return false
			})

			if !ok && item["children"] == nil {
				return item
			} else if item["children"] != nil {
				item["children"] = lo.Map(item["children"].([]any), func(item any, index int) map[string]any {
					return hyperutils.CovertStructToMap(item)
				})
				for idx, child := range item["children"].([]map[string]any) {
					item["children"].([]map[string]any)[idx] = build(child)
				}
			}

			if ok && app.ExposedOptions != nil {
				if page, ok := lo.Find(app.ExposedOptions.Pages, func(v models.SubAppExposedPage) bool {
					return v.Name == item["name"]
				}); ok {
					item["to"] = fmt.Sprintf(
						"%s/srv/subapps/%s%s",
						viper.GetString("base_url"),
						app.Manifest.Name,
						page.Path,
					)
				}
			}

			return item
		}

		nav = append(nav, build(item))
	}

	firmware := strings.Split(c.App().Config().AppName, " ")
	return c.JSON(fiber.Map{
		"debug":            viper.GetBool("debug"),
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
		"pages": pages,
		"nav":   nav,
	})
}
