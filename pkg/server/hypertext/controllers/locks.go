package controllers

import (
	"code.smartsheep.studio/atom/bedrock/pkg/server/datasource/models"
	"code.smartsheep.studio/atom/bedrock/pkg/server/hypertext/hyperutils"
	"code.smartsheep.studio/atom/bedrock/pkg/server/hypertext/middlewares"
	"code.smartsheep.studio/atom/bedrock/pkg/server/services"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"time"
)

type LockController struct {
	db         *gorm.DB
	notify     *services.NotificationService
	gatekeeper *middlewares.AuthMiddleware
}

func NewLockController(db *gorm.DB, notify *services.NotificationService, gatekeeper *middlewares.AuthMiddleware) *LockController {
	return &LockController{db, notify, gatekeeper}
}

func (v *LockController) Map(router *fiber.App) {
	router.Post("/cgi/locks", v.create)

	router.Post(
		"/api/administration/locks",
		v.gatekeeper.Fn(true,
			hyperutils.GenScope("admin:locks"),
			hyperutils.GenPerms("admin.locks.create"),
		),
		v.create,
	)
}

func (v *LockController) create(c *fiber.Ctx) error {
	var req struct {
		Reason    string     `json:"reason"`
		ExpiredAt *time.Time `json:"expired_at"`
		UserID    uint       `json:"user_id"`
	}

	if err := hyperutils.BodyParser(c, &req); err != nil {
		return err
	}

	var user models.User
	if err := v.db.Where("id = ?", req.UserID).First(&user).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}

	item := models.Lock{
		Reason:    req.Reason,
		ExpiredAt: req.ExpiredAt,
		UserID:    &user.ID,
	}

	if err := v.db.Save(&item).Error; err != nil {
		return hyperutils.ErrorParser(err)
	} else {
		if c.Query("quiet", "yes") != "yes" {
			if err := v.notify.SendNotification(&models.Notification{
				Title: "Congratulations!",
				Description: fmt.Sprintf(
					"Congratulations! You have been successfully banned from %s!",
					viper.GetString("general.name"),
				),
				Content: fmt.Sprintf(
					"Congratulations, %s! You have be banned from %s because %s. Now sign in our website to view details!",
					user.Name,
					item.Reason,
					viper.GetString("general.name"),
				),
				Level:       models.NotificationLevelWarning,
				SenderType:  models.NotificationSenderTypeSystem,
				RecipientID: user.ID,
			}); err != nil {
				return hyperutils.ErrorParser(err)
			}
		}

		return c.JSON(err)
	}
}
