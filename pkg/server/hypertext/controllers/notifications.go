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
	"time"
)

type NotificationController struct {
	db            *gorm.DB
	gatekeeper    *middlewares.AuthMiddleware
	notifications *services.NotificationService
}

func NewNotificationController(db *gorm.DB, gatekeeper *middlewares.AuthMiddleware, notifications *services.NotificationService) *NotificationController {
	return &NotificationController{db, gatekeeper, notifications}
}

func (v *NotificationController) Map(router *fiber.App) {
	router.Post("/cgi/notifications", v.send)

	router.Post(
		"/api/administration/notifications",
		v.gatekeeper.Fn(true,
			hyperutils.GenScope("admin:notifications"),
			hyperutils.GenPerms("admin.notifications.create"),
		),
		v.send,
	)

	router.Post(
		"/api/notifications/all/read",
		v.gatekeeper.Fn(true,
			hyperutils.GenScope("read:notifications"),
			hyperutils.GenPerms("notifications.read"),
		),
		v.readAll,
	)
	router.Post(
		"/api/notifications/:notify/read",
		v.gatekeeper.Fn(true,
			hyperutils.GenScope("read:notifications"),
			hyperutils.GenPerms("notifications.read"),
		),
		v.read,
	)
}

func (v *NotificationController) send(c *fiber.Ctx) error {
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

	if err := v.notifications.SendNotification(&item); err != nil {
		return hyperutils.ErrorParser(err)
	} else {
		return c.JSON(item)
	}
}

func (v *NotificationController) read(c *fiber.Ctx) error {
	u := c.Locals("principal").(models.User)

	var notification models.Notification
	if err := v.db.Where("id = ? AND recipient_id = ?", c.Params("notify"), u.ID).First(&notification).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}

	notification.ReadAt = lo.ToPtr(time.Now())

	if err := v.db.Save(&notification).Error; err != nil {
		return hyperutils.ErrorParser(err)
	} else {
		return c.JSON(notification)
	}
}

func (v *NotificationController) readAll(c *fiber.Ctx) error {
	u := c.Locals("principal").(models.User)

	if err := v.db.Where("recipient_id = ?", c.Params("notify"), u.ID).Updates(models.Notification{ReadAt: lo.ToPtr(time.Now())}).Error; err != nil {
		return hyperutils.ErrorParser(err)
	} else {
		return c.SendStatus(fiber.StatusNoContent)
	}
}
