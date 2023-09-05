package controllers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"code.smartsheep.studio/atom/bedrock/pkg/server/datasource/models"
	"code.smartsheep.studio/atom/bedrock/pkg/server/hypertext/hyperutils"
	"code.smartsheep.studio/atom/bedrock/pkg/server/hypertext/middlewares"
	"code.smartsheep.studio/atom/bedrock/pkg/server/services"
	"github.com/samber/lo"
	"github.com/spf13/viper"

	"github.com/gofiber/fiber/v2"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserController struct {
	db         *gorm.DB
	otp        *services.OTPService
	auth       *services.AuthService
	warehouse  *services.StorageService
	gatekeeper *middlewares.AuthMiddleware
}

func NewUserController(db *gorm.DB, otp *services.OTPService, auth *services.AuthService, warehouse *services.StorageService, gatekeeper *middlewares.AuthMiddleware) *UserController {
	ctrl := &UserController{db, otp, auth, warehouse, gatekeeper}
	return ctrl
}

func (v *UserController) Map(router *fiber.App) {
	router.Get(
		"/cgi/users/:user",
		v.info,
	)

	router.Get(
		"/api/administration/users",
		v.gatekeeper.Fn(
			true,
			hyperutils.GenScope("admin:users"),
			hyperutils.GenPerms("admin.users.read"),
		),
		v.list,
	)

	router.Get(
		"/api/users/self",
		v.gatekeeper.Fn(true, hyperutils.GenScope("read:id"), hyperutils.GenPerms()),
		v.self,
	)
	router.Get(
		"/api/users/self/notifications",
		v.gatekeeper.Fn(true, hyperutils.GenScope("read:notifications"), hyperutils.GenPerms()),
		v.selfNotifications,
	)
	router.Get(
		"/api/users/self/locks",
		v.gatekeeper.Fn(true, hyperutils.GenScope("read:locks"), hyperutils.GenPerms()),
		v.selfLocks,
	)

	router.Get(
		"/api/users/self/verify",
		v.gatekeeper.Fn(true, hyperutils.GenScope("verify:id"), hyperutils.GenPerms()),
		v.verify,
	)

	router.Put(
		"/api/users/self",
		v.gatekeeper.Fn(true, hyperutils.GenScope("update:id"), hyperutils.GenPerms()),
		v.update,
	)
	router.Put(
		"/api/users/self/contacts",
		v.gatekeeper.Fn(true, hyperutils.GenScope("update:id.contacts"), hyperutils.GenPerms()),
		v.updateContacts,
	)
	router.Put(
		"/api/users/self/password",
		v.gatekeeper.Fn(true, hyperutils.GenScope("update:id.password"), hyperutils.GenPerms()),
		v.updatePassword,
	)
	router.Put(
		"/api/users/self/personalize",
		v.gatekeeper.Fn(true, hyperutils.GenScope("update:id.personalize"), hyperutils.GenPerms("personalize")),
		v.personalize,
	)
}

func (v *UserController) list(c *fiber.Ctx) error {
	var users []models.User
	if err := v.db.Limit(c.QueryInt("limit", 20)).Offset(c.QueryInt("offset", 0)).Find(&users).Error; err != nil {
		return hyperutils.ErrorParser(err)
	} else {
		return c.JSON(users)
	}
}

func (v *UserController) info(c *fiber.Ctx) error {
	var user models.User
	if err := v.db.Where("id = ?", c.Params("user", "0")).Preload("Contacts").Preload("Groups").First(&user).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}

	var notificationsCount int64
	if err := v.db.Model(&models.Notification{}).Where("recipient_id = ? AND read_at IS NULL", user.ID).Count(&notificationsCount).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}

	m := hyperutils.CovertStructToMap(user)
	m["permissions"], _ = user.GetPermissions()
	m["notifications_count"] = notificationsCount

	return c.JSON(fiber.Map{
		"sub":            user.ID,
		"name":           user.Name,
		"nickname":       user.Nickname,
		"profile":        fmt.Sprintf("%s/explore/users/%s", viper.GetString("general.base_url"), user.Name),
		"email":          user.Contacts[0].Content,
		"email_verified": user.Contacts[0].VerifiedAt != nil,
		"session":        nil,
		"claims":         nil,
		"user":           m,
	})
}

func (v *UserController) self(c *fiber.Ctx) error {
	u := c.Locals("principal").(models.User)

	var user models.User
	if err := v.db.Where("id = ?", u.ID).Preload("Contacts").Preload("Groups").First(&user).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}

	var notificationCount int64
	if err := v.db.Model(&models.Notification{}).Where("recipient_id = ? AND read_at IS NULL", user.ID).Count(&notificationCount).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}

	m := hyperutils.CovertStructToMap(user)
	m["permissions"], _ = user.GetPermissions()
	m["notification_count"] = notificationCount

	claims := c.Locals("principal-claims").(models.UserClaims)

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

func (v *UserController) selfNotifications(c *fiber.Ctx) error {
	u := c.Locals("principal").(models.User)

	tx := v.db.Where("recipient_id = ?", u.ID)
	tx.Order("created_at desc")
	tx.Limit(20)

	if c.Query("only_unread", "yes") == "yes" {
		tx.Where("read_at IS NULL")
	}

	var notifications []models.Notification
	if err := tx.Find(&notifications).Error; err != nil {
		return hyperutils.ErrorParser(err)
	} else {
		return c.JSON(notifications)
	}
}

func (v *UserController) selfLocks(c *fiber.Ctx) error {
	u := c.Locals("principal").(models.User)

	tx := v.db.Where("user_id = ?", u.ID)

	var locks []models.Lock
	if err := tx.Find(&locks).Error; err != nil {
		return hyperutils.ErrorParser(err)
	} else {
		return c.JSON(locks)
	}
}

func (v *UserController) verify(c *fiber.Ctx) error {
	u := c.Locals("principal").(models.User)

	code := c.Query("code")
	if len(code) > 0 {
		otp, err := v.otp.LookupOTP(u, c.Query("code"))
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		if err := v.otp.ApplyOTP(otp); err != nil {
			return hyperutils.ErrorParser(err)
		} else {
			return c.SendString("Successfully verified your account and contact!")
		}
	}

	id := c.QueryInt("id", 0)
	if id <= 0 {
		return fiber.NewError(fiber.StatusBadRequest, "You need provide a valid contact id.")
	}

	var contact models.UserContact
	if err := v.db.Where("id = ? AND user_id = ?", id, u.ID).First(&contact).Error; err != nil {
		return hyperutils.ErrorParser(err)
	} else if contact.VerifiedAt != nil {
		return fiber.NewError(fiber.StatusBadRequest, "You need provide a unverified contact id.")
	}

	exp, _ := time.ParseDuration("30m")
	otp, err := v.otp.NewOTP(u, models.OneTimeVerifyContactCode, models.OTPPayload{
		Target:    strconv.Itoa(int(contact.ID)),
		IpAddress: c.IP(),
	}, &exp)
	if err != nil {
		return fiber.NewError(fiber.StatusForbidden, err.Error())
	} else {
		if err := v.otp.SendMail(otp); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		return c.SendString("Verify email has been sent.")
	}
}

func (v *UserController) update(c *fiber.Ctx) error {
	var req struct {
		Name        string `json:"name" validate:"required"`
		Nickname    string `json:"nickname" validate:"required"`
		Description string `json:"description"`
	}

	if err := hyperutils.BodyParser(c, &req); err != nil {
		return err
	}

	user := c.Locals("principal").(models.User)
	user.Name = strings.ToLower(req.Name)
	user.Nickname = req.Nickname
	user.Description = req.Description

	if err := v.db.Save(&user).Error; err != nil {
		return hyperutils.ErrorParser(err)
	} else {
		return c.JSON(user)
	}
}

func (v *UserController) updateContacts(c *fiber.Ctx) error {
	u := c.Locals("principal").(models.User)

	var req struct {
		Contacts []struct {
			Type        string `json:"type" validate:"required"`
			Content     string `json:"content" validate:"required"`
			Description string `json:"description"`
		} `json:"contacts" validate:"required"`
	}

	if err := hyperutils.BodyParser(c, &req); err != nil {
		return err
	}

	tx := v.db.Begin()

	if err := tx.Unscoped().Delete(&u.Contacts).Error; err != nil {
		tx.Rollback()
		return hyperutils.ErrorParser(err)
	}

	var contacts []models.UserContact
	for _, item := range req.Contacts {
		var old *models.UserContact
		for _, record := range u.Contacts {
			if record.Content == item.Content {
				old = &record
				break
			}
		}
		contacts = append(contacts, models.UserContact{
			Type:        item.Type,
			Content:     item.Content,
			Description: item.Description,
			VerifiedAt:  lo.Ternary(old != nil, old.VerifiedAt, nil),
		})
	}

	u.Contacts = contacts

	if err := tx.Save(&u).Error; err != nil {
		tx.Rollback()
		return hyperutils.ErrorParser(err)
	}

	if err := tx.Commit().Error; err != nil {
		return hyperutils.ErrorParser(err)
	} else {
		return c.JSON(u)
	}
}

func (v *UserController) updatePassword(c *fiber.Ctx) error {
	u := c.Locals("principal").(models.User)

	var req struct {
		OldPassword string `json:"old_password" validate:"required"`
		NewPassword string `json:"new_password" validate:"required"`
	}

	if err := hyperutils.BodyParser(c, &req); err != nil {
		return err
	}

	var encrypted []byte
	if _, err := v.auth.AuthUser(u.Name, req.OldPassword); err != nil {
		return fiber.NewError(fiber.StatusForbidden, "invalid old password")
	} else if encrypted, err = bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to process your password: %s", err.Error()))
	}

	u.Password = string(encrypted)

	if err := v.db.Save(&u).Error; err != nil {
		return hyperutils.ErrorParser(err)
	} else {
		return c.JSON(u)
	}
}

func (v *UserController) personalize(c *fiber.Ctx) error {
	u := c.Locals("principal").(models.User)

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
		v.db.Delete(&models.StorageFile{UserID: &u.ID, Type: models.StorageFileAvatarType})

		if f, err := v.warehouse.SaveFile2User(c, avatar, u, models.StorageFileAvatarType); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		} else {
			// Update url
			u.AvatarUrl = f.GetURL()
			v.db.Save(&u)

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
		v.db.Delete(&models.StorageFile{UserID: &u.ID, Type: models.StorageFileBannerType})

		if f, err := v.warehouse.SaveFile2User(c, banner, u, models.StorageFileBannerType); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		} else {
			// Update url
			u.BannerUrl = f.GetURL()
			v.db.Save(&u)

			return c.SendStatus(fiber.StatusOK)
		}
	default:
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("unexpected field: %s", field))
	}
}
