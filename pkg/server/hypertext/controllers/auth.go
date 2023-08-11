package controllers

import (
	models2 "code.smartsheep.studio/atom/bedrock/pkg/server/datasource/models"
	hyperutils2 "code.smartsheep.studio/atom/bedrock/pkg/server/hypertext/hyperutils"
	"code.smartsheep.studio/atom/bedrock/pkg/server/hypertext/middlewares"
	services2 "code.smartsheep.studio/atom/bedrock/pkg/server/services"
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
	auth       *services2.AuthService
	users      *services2.UserService
	gatekeeper *middlewares.AuthMiddleware
}

func NewAuthController(db *gorm.DB, auth *services2.AuthService, users *services2.UserService, gatekeeper *middlewares.AuthMiddleware) *AuthController {
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
		ctrl.gatekeeper.Fn(true, hyperutils2.GenScope("read:id.sessions"), hyperutils2.GenPerms()),
		ctrl.securitySessions,
	)
	router.Delete(
		"/api/auth/sessions",
		ctrl.gatekeeper.Fn(true, hyperutils2.GenScope("delete:id.sessions"), hyperutils2.GenPerms()),
		ctrl.securityTerminate,
	)
}

func (ctrl *AuthController) signin(c *fiber.Ctx) error {
	var req struct {
		ID       string `json:"id" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	if err := hyperutils2.BodyParser(c, &req); err != nil {
		return err
	}

	user, err := ctrl.auth.AuthUser(req.ID, req.Password)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid username or password")
	}

	exp := time.Now().Add(viper.GetDuration("security.sessions_alive_duration"))
	session := models2.UserSession{
		IpAddress: c.IP(),
		Type:      models2.UserSessionTypeAuth,
		UserID:    user.ID,
		ExpiredAt: &exp,
		Scope:     datatypes.NewJSONSlice([]string{"*"}),
		Available: true,
	}

	if err := ctrl.db.Save(&session).Error; err != nil {
		return hyperutils2.ErrorParser(err)
	}

	_, token, err := ctrl.auth.NewJwt(session, models2.UserClaimsTypeAccess)
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

	if err := hyperutils2.BodyParser(c, &req); err != nil {
		return err
	}

	contact := models2.UserContact{
		Name:        "Primary Contact",
		Description: fmt.Sprintf("%s's Primary Contact", req.Name),
		Content:     req.Contact,
		IsPrimary:   true,
	}

	if ok, _ := regexp.Match("^[\\w-\\.]+@([\\w-]+\\.)+[\\w-]{2,4}$", []byte(req.Contact)); ok {
		contact.Type = models2.UserContactTypeEmail
	} else if ok, _ := regexp.Match("\\+?\\d{1,4}?[-.\\s]?\\(?\\d{1,3}?\\)?[-.\\s]?\\d{1,4}[-.\\s]?\\d{1,4}[-.\\s]?\\d{1,9}\n", []byte(req.Contact)); ok {
		contact.Type = models2.UserContactTypePhone
	} else {
		return fiber.NewError(fiber.StatusBadRequest, "invalid contact type")
	}

	item := models2.User{
		Name:        req.Name,
		Nickname:    req.Nickname,
		Description: "The man is too lazy to write anything.",
		Password:    req.Password,
		Contacts:    []models2.UserContact{contact},
	}

	if err := ctrl.users.NewUser(&item); err != nil {
		return hyperutils2.ErrorParser(err)
	} else {
		return c.Status(fiber.StatusCreated).JSON(item)
	}
}

func (ctrl *AuthController) securitySessions(c *fiber.Ctx) error {
	u := c.Locals("principal").(models2.User)

	var sessions []models2.UserSession
	if err := ctrl.db.Where("user_id = ?", u.ID).Find(&sessions).Error; err != nil {
		return hyperutils2.ErrorParser(err)
	}

	// Load mentioned clients
	clients := map[uint]*models2.OauthClient{}
	for _, session := range sessions {
		if session.ClientID != nil {
			if clients[*session.ClientID] == nil {
				var client models2.OauthClient
				if err := ctrl.db.Where("id = ?", session.ClientID).First(&client).Error; err != nil {
					return hyperutils2.ErrorParser(err)
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
	u := c.Locals("principal").(models2.User)

	id := strconv.Itoa(int(c.Locals("principal-session").(models2.UserSession).ID))
	if len(c.Query("id")) > 0 {
		id = c.Query("id")
	}

	if err := ctrl.db.Where("user_id = ? AND id = ?", u.ID, id).Delete(&models2.UserSession{}).Error; err != nil {
		return hyperutils2.ErrorParser(err)
	} else {
		return c.SendStatus(fiber.StatusNoContent)
	}
}
