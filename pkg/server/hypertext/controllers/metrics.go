package controllers

import (
	"code.smartsheep.studio/atom/bedrock/pkg/server/datasource/models"
	"code.smartsheep.studio/atom/bedrock/pkg/server/hypertext/hyperutils"
	"code.smartsheep.studio/atom/bedrock/pkg/server/hypertext/middlewares"
	"code.smartsheep.studio/atom/bedrock/pkg/server/services"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"time"
)

type MetricsController struct {
	db         *gorm.DB
	metrics    *services.MetricsService
	gatekeeper *middlewares.AuthMiddleware
}

func NewMetricsController(db *gorm.DB, metrics *services.MetricsService, gatekeeper *middlewares.AuthMiddleware) *MetricsController {
	return &MetricsController{db, metrics, gatekeeper}
}

func (v *MetricsController) Map(router *fiber.App) {
	router.Get(
		"/api/metrics",
		v.gatekeeper.Fn(
			true,
			hyperutils.GenScope("read:metrics"),
			hyperutils.GenPerms("admin.metrics.read"),
		),
		v.overview,
	)
}

func (v *MetricsController) overview(c *fiber.Ctx) error {
	countRecord := func(m any) int64 {
		var count int64
		if err := v.db.Model(m).Count(&count).Error; err != nil {
			return -1
		} else {
			return count
		}
	}

	return c.JSON(fiber.Map{
		"uptime": time.Since(v.metrics.StartAt).Milliseconds(),
		"resources": fiber.Map{
			"users":         countRecord(&models.User{}),
			"sessions":      countRecord(&models.UserSession{}),
			"contacts":      countRecord(&models.UserContact{}),
			"notifications": countRecord(&models.Notification{}),
		},
	})
}
