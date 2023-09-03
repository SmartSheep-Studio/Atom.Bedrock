package controllers

import (
	"code.smartsheep.studio/atom/bedrock/pkg/server/datasource/models"
	"code.smartsheep.studio/atom/bedrock/pkg/server/hypertext/hyperutils"
	"code.smartsheep.studio/atom/bedrock/pkg/server/hypertext/middlewares"
	"code.smartsheep.studio/atom/bedrock/pkg/server/services"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type NotificationController struct {
	db         *gorm.DB
	gatekeeper *middlewares.AuthMiddleware
	notifications *services.NotificationService
}

func NewNotificationController(db *gorm.DB, gatekeeper *middlewares.AuthMiddleware, notifications *services.NotificationService) *NotificationController {
	return &NotificationController{db, gatekeeper, notifications}
}

func (v *NotificationController) Map(router *fiber.App) {
	router.Post("/cgi/notifications", v.SendNotification)

	router.Post(
		"/api/administration/notifications",
		v.gatekeeper.Fn(true,
			hyperutils.GenScope("admin:notifications"),
			hyperutils.GenPerms("admin.notifications.create"),
		),
		v.SendNotification,
	)
}

func (v *NotificationController) SendNotification(c *fiber.Ctx) error {
	var req struct {
		Title       string                    `json:"title" validate:"required"`
		Description string                    `json:"description" validate:"required"`
		Content     string                    `json:"content" validate:"required"`
		Level       string                    `json:"level" validate:"required"`
		Links       []models.NotificationLink `json:"links"`
		RecipientId uint                      `json:"recipient_id" validate:"required"`
		SenderID    *uint                     `json:"sender_id"`
	}

	if err := hyperutils.BodyParser(c, &req); err != nil {
		return err
	}

	item := models.Notification{
		Title:       req.Title,
		Description: req.Description,
		Content:     req.Content,
		Level:       req.Level,
		Links:       datatypes.NewJSONSlice(req.Links),
		RecipientID: req.RecipientId,
		SenderType:  lo.Ternary(req.SenderID == nil, models.NotificationSenderTypeSystem, models.NotificationSenderTypeUser),
		SenderID:    req.SenderID,
	}
	
	if err := v.notifications.SendNotification(item).Error; err != nil {
		return hyperutils.ErrorParser(err)
	} else {
		return c.JSON(item)
	}
}
