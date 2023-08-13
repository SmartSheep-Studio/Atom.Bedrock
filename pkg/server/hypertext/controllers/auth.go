package controllers

import (
	models "code.smartsheep.studio/atom/bedrock/pkg/server/datasource/models"
	hyperutils "code.smartsheep.studio/atom/bedrock/pkg/server/hypertext/hyperutils"
	"code.smartsheep.studio/atom/bedrock/pkg/server/hypertext/middlewares"
	services "code.smartsheep.studio/atom/bedrock/pkg/server/services"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/spf13/viper"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type AuthController struct {
	db         *gorm.DB
	auth       *services.AuthService
	users      *services.UserService
	gatekeeper *middlewares.AuthMiddleware
}

func NewAuthController(db *gorm.DB, auth *services.AuthService, users *services.UserService, gatekeeper *middlewares.AuthMiddleware) *AuthController {
	ctrl := &AuthController{db, auth, users, gatekeeper}
	return ctrl
}

func (ctrl *AuthController) Map(router *fiber.App) {
	router.Post(
		"/api/auth/sign-in",
		ctrl.signin,
	)
	router.Post(
		"/api/auth/sign-up",
		ctrl.signup,
	)
	router.Get("/api/auth/sessions",
		ctrl.gatekeeper.Fn(true, hyperutils.GenScope("read:id.sessions"), hyperutils.GenPerms()),
		ctrl.securitySessions,
	)
	router.Delete(
		"/api/auth/sessions",
		ctrl.gatekeeper.Fn(true, hyperutils.GenScope("delete:id.sessions"), hyperutils.GenPerms()),
		ctrl.securityTerminate,
	)
}

func (ctrl *AuthController) signin(c *fiber.Ctx) error {
	var req struct {
		ID       string `json:"id" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	if err := hyperutils.BodyParser(c, &req); err != nil {
		return err
	}

	user, err := ctrl.auth.AuthUser(req.ID, req.Password)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid username or password")
	} else {
		for _, lock := range user.Locks {
			if lock.ExpiredAt == nil || lock.ExpiredAt.Unix() >= time.Now().Unix() {
				return fiber.NewError(
					fiber.StatusForbidden,
					fmt.Sprintf("your account has been locked, reason: %s", lock.Reason),
				)
			}
		}
	}

	exp := time.Now().Add(viper.GetDuration("security.sessions_alive_duration"))
	session := models.UserSession{
		IpAddress: c.IP(),
		Type:      models.UserSessionTypeAuth,
		UserID:    user.ID,
		ExpiredAt: &exp,
		Scope:     datatypes.NewJSONSlice([]string{"*"}),
		Available: true,
	}

	if err := ctrl.db.Save(&session).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}

	_, token, err := ctrl.auth.NewJwt(session, models.UserClaimsTypeAccess)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to encode jwt: %s", err.Error()))
	} else {
		c.Cookie(&fiber.Cookie{
			Path:     "/",
			Name:     "authorization",
			Value:    token,
			MaxAge:   int(viper.GetDuration("security.sessions_alive_duration").Seconds()),
			Domain:   viper.GetString("security.cookies_domain"),
			SameSite: "None",
			Secure:   true,
		})

		return c.JSON(fiber.Map{
			"user":  user,
			"token": token,
		})
	}
}

func (ctrl *AuthController) signup(c *fiber.Ctx) error {
	var req struct {
		Name     string `json:"name" validate:"required,min=4,max=16"`
		Nickname string `json:"nickname" validate:"required,min=4,max=16"`
		Password string `json:"password" validate:"required,min=8,max=32"`
		Contact  string `json:"contact" validate:"required"`
	}

	if err := hyperutils.BodyParser(c, &req); err != nil {
		return err
	}

	contact := models.UserContact{
		Name:        "Primary Contact",
		Description: fmt.Sprintf("%s's Primary Contact", req.Name),
		Content:     req.Contact,
		IsPrimary:   true,
	}

	if ok, _ := regexp.Match("^[\\w-\\.]+@([\\w-]+\\.)+[\\w-]{2,4}$", []byte(req.Contact)); ok {
		contact.Type = models.UserContactTypeEmail
	} else if ok, _ := regexp.Match("\\+?\\d{1,4}?[-.\\s]?\\(?\\d{1,3}?\\)?[-.\\s]?\\d{1,4}[-.\\s]?\\d{1,4}[-.\\s]?\\d{1,9}\n", []byte(req.Contact)); ok {
		contact.Type = models.UserContactTypePhone
	} else {
		return fiber.NewError(fiber.StatusBadRequest, "invalid contact type")
	}

	item := models.User{
		Name:        req.Name,
		Nickname:    req.Nickname,
		Description: "The man is too lazy to write anything.",
		Password:    req.Password,
		Contacts:    []models.UserContact{contact},
	}

	if err := ctrl.users.NewUser(&item); err != nil {
		return hyperutils.ErrorParser(err)
	} else {
		return c.Status(fiber.StatusCreated).JSON(item)
	}
}

func (ctrl *AuthController) securitySessions(c *fiber.Ctx) error {
	u := c.Locals("principal").(models.User)

	var sessions []models.UserSession
	if err := ctrl.db.Where("user_id = ?", u.ID).Find(&sessions).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}

	// Load mentioned clients
	clients := map[uint]*models.OauthClient{}
	for _, session := range sessions {
		if session.ClientID != nil {
			if clients[*session.ClientID] == nil {
				var client models.OauthClient
				if err := ctrl.db.Where("id = ?", session.ClientID).First(&client).Error; err != nil {
					return hyperutils.ErrorParser(err)
				} else {
					clients[*session.ClientID] = &client
				}
			}
		}
	}

	return c.JSON(fiber.Map{
		"clients":  clients,
		"sessions": sessions,
	})
}

func (ctrl *AuthController) securityTerminate(c *fiber.Ctx) error {
	u := c.Locals("principal").(models.User)

	id := strconv.Itoa(int(c.Locals("principal-session").(models.UserSession).ID))
	if len(c.Query("id")) > 0 {
		id = c.Query("id")
	}

	if err := ctrl.db.Where("user_id = ? AND id = ?", u.ID, id).Delete(&models.UserSession{}).Error; err != nil {
		return hyperutils.ErrorParser(err)
	} else {
		return c.SendStatus(fiber.StatusNoContent)
	}
}
