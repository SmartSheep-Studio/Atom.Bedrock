package controllers

import (
	models2 "code.smartsheep.studio/atom/bedrock/pkg/server/datasource/models"
	hyperutils2 "code.smartsheep.studio/atom/bedrock/pkg/server/hypertext/hyperutils"
	"code.smartsheep.studio/atom/bedrock/pkg/server/hypertext/middlewares"
	services2 "code.smartsheep.studio/atom/bedrock/pkg/server/services"
	"fmt"
	"github.com/samber/lo"
	"github.com/spf13/viper"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserController struct {
	db         *gorm.DB
	otp        *services2.OTPService
	auth       *services2.AuthService
	warehouse  *services2.StorageService
	gatekeeper *middlewares.AuthMiddleware
}

func NewUserController(db *gorm.DB, otp *services2.OTPService, auth *services2.AuthService, warehouse *services2.StorageService, gatekeeper *middlewares.AuthMiddleware) *UserController {
	ctrl := &UserController{db, otp, auth, warehouse, gatekeeper}
	return ctrl
}

func (ctrl *UserController) Map(router *fiber.App) {
	router.Get(
		"/api/users/self",
		ctrl.gatekeeper.Fn(true, hyperutils2.GenScope("read:id"), hyperutils2.GenPerms()),
		ctrl.self,
	)
	router.Get(
		"/api/users/self/notifications",
		ctrl.gatekeeper.Fn(true, hyperutils2.GenScope("read:notifications"), hyperutils2.GenPerms()),
		ctrl.selfNotifications,
	)

	router.Get(
		"/api/users/self/verify",
		ctrl.gatekeeper.Fn(true, hyperutils2.GenScope("verify:id"), hyperutils2.GenPerms()),
		ctrl.verify,
	)

	router.Put(
		"/api/users/self",
		ctrl.gatekeeper.Fn(true, hyperutils2.GenScope("update:id"), hyperutils2.GenPerms()),
		ctrl.update,
	)
	router.Put(
		"/api/users/self/contacts",
		ctrl.gatekeeper.Fn(true, hyperutils2.GenScope("update:id.contacts"), hyperutils2.GenPerms()),
		ctrl.updateContacts,
	)
	router.Put(
		"/api/users/self/password",
		ctrl.gatekeeper.Fn(true, hyperutils2.GenScope("update:id.password"), hyperutils2.GenPerms()),
		ctrl.updatePassword,
	)
	router.Put(
		"/api/users/self/personalize",
		ctrl.gatekeeper.Fn(true, hyperutils2.GenScope("update:id.personalize"), hyperutils2.GenPerms("personalize")),
		ctrl.personalize,
	)
}

func (ctrl *UserController) self(c *fiber.Ctx) error {
	u := c.Locals("principal").(models2.User)

	var user models2.User
	if err := ctrl.db.Where("id = ?", u.ID).Preload("Contacts").First(&user).Error; err != nil {
		return hyperutils2.ErrorParser(err)
	}

	m := hyperutils2.CovertStructToMap(user)
	m["permissions"], _ = user.GetPermissions()

	claims := c.Locals("principal-claims").(models2.UserClaims)

	return c.JSON(fiber.Map{
		"sub":            claims.Subject,
		"name":           user.Name,
		"nickname":       user.Nickname,
		"profile":        fmt.Sprintf("%s/explore/users/%s", viper.GetString("general.base_url"), user.Name),
		"email":          user.Contacts[0].Content,
		"email_verified": user.Contacts[0].VerifiedAt != nil,
		"claims":         c.Locals("principal-claims"),
		"session":        c.Locals("principal-session"),
		"user":           m,
	})
}

func (ctrl *UserController) selfNotifications(c *fiber.Ctx) error {
	u := c.Locals("principal").(models2.User)

	tx := ctrl.db.Where("recipient_id = ?", u.ID)
	if c.Query("only_unread", "yes") == "yes" {
		tx.Where("read_at IS NULL")
	}

	var notifications []models2.Notification
	if err := tx.Find(&notifications).Error; err != nil {
		return hyperutils2.ErrorParser(err)
	} else {
		if c.Query("update_state", "yes") == "yes" {
			ctrl.db.Model(models2.Notification{}).Where("recipient_id = ? AND read_at IS NULL", u.ID).Updates(models2.Notification{
				ReadAt: lo.ToPtr(time.Now()),
			})
		}

		return c.JSON(notifications)
	}
}

func (ctrl *UserController) verify(c *fiber.Ctx) error {
	u := c.Locals("principal").(models2.User)

	code := c.Query("code")
	if len(code) > 0 {
		otp, err := ctrl.otp.LookupOTP(u, c.Query("code"))
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		if err := ctrl.otp.ApplyOTP(otp); err != nil {
			return hyperutils2.ErrorParser(err)
		} else {
			return c.SendString("Successfully verified your account and contact!")
		}
	}

	id := c.QueryInt("id", 0)
	if id <= 0 {
		return fiber.NewError(fiber.StatusBadRequest, "You need provide a valid contact id.")
	}

	var contact models2.UserContact
	if err := ctrl.db.Where("id = ? AND user_id = ?", id, u.ID).First(&contact).Error; err != nil {
		return hyperutils2.ErrorParser(err)
	} else if contact.VerifiedAt != nil {
		return fiber.NewError(fiber.StatusBadRequest, "You need provide a unverified contact id.")
	}

	exp, _ := time.ParseDuration("30m")
	otp, err := ctrl.otp.NewOTP(u, models2.OneTimeVerifyContactCode, models2.OTPPayload{
		Target:    strconv.Itoa(int(contact.ID)),
		IpAddress: c.IP(),
	}, &exp)
	if err != nil {
		return fiber.NewError(fiber.StatusForbidden, err.Error())
	} else {
		if err := ctrl.otp.SendMail(otp); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		return c.SendString("Verify email has been sent.")
	}
}

func (ctrl *UserController) update(c *fiber.Ctx) error {
	var req struct {
		Name        string `json:"name" validate:"required"`
		Nickname    string `json:"nickname" validate:"required"`
		Description string `json:"description"`
	}

	if err := hyperutils2.BodyParser(c, &req); err != nil {
		return err
	}

	user := c.Locals("principal").(models2.User)
	user.Name = strings.ToLower(req.Name)
	user.Nickname = req.Nickname
	user.Description = req.Description

	if err := ctrl.db.Save(&user).Error; err != nil {
		return hyperutils2.ErrorParser(err)
	} else {
		return c.JSON(user)
	}
}

func (ctrl *UserController) updateContacts(c *fiber.Ctx) error {
	u := c.Locals("principal").(models2.User)

	var req struct {
		Contacts []struct {
			Type        string `json:"type" validate:"required"`
			Content     string `json:"content" validate:"required"`
			Description string `json:"description"`
		} `json:"contacts" validate:"required"`
	}

	if err := hyperutils2.BodyParser(c, &req); err != nil {
		return err
	}

	tx := ctrl.db.Begin()

	if err := tx.Unscoped().Delete(&u.Contacts).Error; err != nil {
		tx.Rollback()
		return hyperutils2.ErrorParser(err)
	}

	var contacts []models2.UserContact
	for _, item := range req.Contacts {
		var old *models2.UserContact
		for _, record := range u.Contacts {
			if record.Content == item.Content {
				old = &record
				break
			}
		}
		contacts = append(contacts, models2.UserContact{
			Type:        item.Type,
			Content:     item.Content,
			Description: item.Description,
			VerifiedAt:  lo.Ternary(old != nil, old.VerifiedAt, nil),
		})
	}

	u.Contacts = contacts

	if err := tx.Save(&u).Error; err != nil {
		tx.Rollback()
		return hyperutils2.ErrorParser(err)
	}

	if err := tx.Commit().Error; err != nil {
		return hyperutils2.ErrorParser(err)
	} else {
		return c.JSON(u)
	}
}

func (ctrl *UserController) updatePassword(c *fiber.Ctx) error {
	u := c.Locals("principal").(models2.User)

	var req struct {
		OldPassword string `json:"old_password" validate:"required"`
		NewPassword string `json:"new_password" validate:"required"`
	}

	if err := hyperutils2.BodyParser(c, &req); err != nil {
		return err
	}

	var encrypted []byte
	if _, err := ctrl.auth.AuthUser(u.Name, req.OldPassword); err != nil {
		return fiber.NewError(fiber.StatusForbidden, "invalid old password")
	} else if encrypted, err = bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to process your password: %s", err.Error()))
	}

	u.Password = string(encrypted)

	if err := ctrl.db.Save(&u).Error; err != nil {
		return hyperutils2.ErrorParser(err)
	} else {
		return c.JSON(u)
	}
}

func (ctrl *UserController) personalize(c *fiber.Ctx) error {
	u := c.Locals("principal").(models2.User)

	field := c.Query("field", "none")
	switch field {
	case "avatar":
		avatar, err := c.FormFile("avatar")
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		} else if !strings.HasPrefix(avatar.Header.Get("Content-Type"), "image") {
			return fiber.NewError(fiber.StatusBadRequest, "banner image only accept images")
		}

		// Clean up old avatars
		ctrl.db.Delete(&models2.StorageFile{UserID: &u.ID, Type: models2.StorageFileAvatarType})

		if f, err := ctrl.warehouse.SaveFile2User(c, avatar, u, models2.StorageFileAvatarType); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		} else {
			// Update url
			u.AvatarUrl = f.GetURL()
			ctrl.db.Save(&u)

			return c.SendStatus(fiber.StatusOK)
		}
	case "banner":
		banner, err := c.FormFile("banner")
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		} else if !strings.HasPrefix(banner.Header.Get("Content-Type"), "image") {
			return fiber.NewError(fiber.StatusBadRequest, "banner image only accept images")
		}

		// Clean up old banners
		ctrl.db.Delete(&models2.StorageFile{UserID: &u.ID, Type: models2.StorageFileBannerType})

		if f, err := ctrl.warehouse.SaveFile2User(c, banner, u, models2.StorageFileBannerType); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		} else {
			// Update url
			u.BannerUrl = f.GetURL()
			ctrl.db.Save(&u)

			return c.SendStatus(fiber.StatusOK)
		}
	default:
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("unexpected field: %s", field))
	}
}
