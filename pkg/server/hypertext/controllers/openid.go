package controllers

import (
	models "code.smartsheep.studio/atom/bedrock/pkg/server/datasource/models"
	hyperutils "code.smartsheep.studio/atom/bedrock/pkg/server/hypertext/hyperutils"
	"code.smartsheep.studio/atom/bedrock/pkg/server/hypertext/middlewares"
	"code.smartsheep.studio/atom/bedrock/pkg/server/services"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/google/uuid"
	"github.com/spf13/viper"
	"golang.org/x/exp/slices"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type OpenIDController struct {
	db         *gorm.DB
	auth       *services.AuthService
	gatekeeper *middlewares.AuthMiddleware
}

func NewOpenIDController(db *gorm.DB, auth *services.AuthService, gatekeeper *middlewares.AuthMiddleware) *OpenIDController {
	ctrl := &OpenIDController{db, auth, gatekeeper}
	return ctrl
}

func (ctrl *OpenIDController) Map(router *fiber.App) {
	router.Get(
		"/api/auth/openid/connect",
		ctrl.gatekeeper.Fn(false, hyperutils.GenScope(), hyperutils.GenPerms()),
		ctrl.connect,
	)
	router.Post(
		"/api/auth/openid/connect",
		ctrl.gatekeeper.Fn(false, hyperutils.GenScope(), hyperutils.GenPerms()),
		ctrl.approve,
	)
	router.Post(
		"/api/auth/openid/exchange",
		ctrl.exchange,
	)
}

func (ctrl *OpenIDController) connect(c *fiber.Ctx) error {
	ok := c.Locals("principal-ok").(bool)

	id := c.Query("client_id")
	redirect := c.Query("redirect_uri")

	var client models.OauthClient
	if err := ctrl.db.Where("slug = ?", id).First(&client).Error; err != nil {
		return hyperutils.ErrorParser(err)
	} else if !client.IsDeveloping && !slices.Contains(client.Callbacks, strings.Split(redirect, "?")[0]) {
		return fiber.NewError(fiber.StatusForbidden, "invalid request url")
	}

	if !ok {
		// Handle unauthorized
		// (Only return client information)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"client":  client,
			"session": nil,
		})
	}

	u := c.Locals("principal").(models.User)

	var session models.UserSession
	if err := ctrl.db.Where("user_id = ? AND client_id = ? AND ip_address = ?", u.ID, client.ID, c.IP()).First(&session); err == nil {
		if session.ExpiredAt.Unix() < time.Now().Unix() {
			return c.JSON(fiber.Map{
				"client":  client,
				"session": nil,
			})
		}

		session.Code = strings.Replace(uuid.New().String(), "-", "", -1)
		ctrl.db.Save(&session)

		return c.JSON(fiber.Map{
			"client":  client,
			"session": session,
		})
	}

	return c.JSON(fiber.Map{
		"client":  client,
		"session": nil,
	})
}

func (ctrl *OpenIDController) approve(c *fiber.Ctx) error {
	var u models.User
	var req struct {
		ID       string `json:"id" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	err := hyperutils.BodyParser(c, &req)
	if !c.Locals("principal-ok").(bool) && err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	} else if !c.Locals("principal-ok").(bool) {
		if u, err = ctrl.auth.AuthUser(req.ID, req.Password); err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
		}
	} else {
		u = c.Locals("principal").(models.User)
	}

	id := c.Query("client_id")
	response := c.Query("response_type")
	redirect := c.Query("redirect_uri")
	scope := c.Query("scope")
	if len(scope) <= 0 {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request params")
	}

	var client models.OauthClient
	if err := ctrl.db.Where("slug = ?", id).First(&client).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}

	if response == "code" {
		exp := time.Now().Add(7 * 24 * time.Hour)
		code := strings.Replace(uuid.New().String(), "-", "", -1)
		session := models.UserSession{
			Type:      models.UserSessionTypeOauth,
			Code:      code,
			Scope:     datatypes.NewJSONSlice(strings.Split(scope, " ")),
			IpAddress: c.IP(),
			Location:  "Unknown",
			Available: true,
			ExpiredAt: &exp,
			ClientID:  &client.ID,
			UserID:    u.ID,
		}

		if err := ctrl.db.Save(&session).Error; err != nil {
			return hyperutils.ErrorParser(err)
		} else {
			return c.JSON(fiber.Map{
				"session":      session,
				"redirect_uri": redirect,
			})
		}
	} else if response == "token" {
		// OAuth Implicit Mode
		exp := time.Now().Add(24 * 14 * time.Hour)
		session := models.UserSession{
			ExpiredAt: &exp,
			Type:      models.UserSessionTypeOauth,
			Scope:     datatypes.NewJSONSlice(strings.Split(scope, " ")),
			IpAddress: c.IP(),
			Location:  "Unknown",
			Available: true,
			ClientID:  &client.ID,
			UserID:    u.ID,
		}

		access, accessToken, err := ctrl.auth.NewJwt(session, models.UserClaimsTypeAccess)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		refresh, refreshToken, err := ctrl.auth.NewJwt(session, models.UserClaimsTypeRefresh)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		session.Access = access.ID
		session.Refresh = refresh.ID

		if err := ctrl.db.Save(&session).Error; err != nil {
			return hyperutils.ErrorParser(err)
		} else {
			return c.JSON(fiber.Map{
				"access_token":  accessToken,
				"refresh_token": refreshToken,
				"session":       session,
				"redirect_uri":  redirect,
			})
		}
	}

	return fiber.NewError(fiber.StatusBadRequest, "unsupported response type")
}

func (ctrl *OpenIDController) exchange(c *fiber.Ctx) error {
	var req struct {
		Code         string `json:"code" form:"code"`
		GrantType    string `json:"grant_type" form:"grant_type" validate:"required"`
		Redirect     string `json:"redirect_uri" form:"redirect_uri"`
		ID           string `json:"username" form:"username"`
		Password     string `json:"password" form:"password"`
		Scope        string `json:"scope" form:"scope"`
		ClientID     string `json:"client_id" form:"client_id"`
		ClientSecret string `json:"client_secret" form:"client_secret"`
	}

	if err := hyperutils.BodyParser(c, &req); err != nil {
		return err
	} else if !slices.Contains([]string{"authorization_code", "refresh_token", "password"}, req.GrantType) {
		return fiber.NewError(fiber.StatusBadRequest, "unsupported grant type")
	}

	var err error
	var user models.User
	var client models.OauthClient
	var session models.UserSession
	if req.GrantType == "refresh_token" {
		// OAuth Refresh Mode
		var pl struct {
			GrantType    string `json:"grant_type" validate:"required"`
			RefreshToken string `json:"refresh_token" validate:"required"`
		}

		if err := c.BodyParser(&pl); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		if claims, err := ctrl.auth.ReadJwt(pl.RefreshToken); err == nil {
			if claims.Type != models.UserClaimsTypeRefresh {
				return fiber.NewError(fiber.StatusForbidden, "invalid token: type check failed")
			}

			id, _ := strconv.Atoi(claims.Subject)
			if err = ctrl.db.Where("id = ?", id).First(&user).Error; err != nil {
				return hyperutils.ErrorParser(err)
			}
			if err = ctrl.db.Where("id = ?", *claims.ClientID).First(&client).Error; err != nil {
				return hyperutils.ErrorParser(err)
			}
			if err = ctrl.db.Where("refresh = ?", claims.ID).First(&session).Error; err != nil {
				return hyperutils.ErrorParser(err)
			}
		} else {
			return fiber.NewError(fiber.StatusForbidden, "invalid token: parse failed")
		}
	} else if req.GrantType == "authorization_code" {
		// OAuth Authorization Code Mode
		if err = ctrl.db.Where("code = ?", req.Code).First(&session).Error; err != nil {
			return hyperutils.ErrorParser(err)
		} else {
			ctrl.db.Where("id = ?", session.UserID).First(&user)
		}

		if err = ctrl.db.Where("slug = ?", req.ClientID).First(&client).Error; err != nil {
			return hyperutils.ErrorParser(err)
		} else if client.Secret != req.ClientSecret {
			return fiber.NewError(fiber.StatusForbidden, "invalid client secret")
		}
	} else if req.GrantType == "password" {
		// OAuth Password Mode
		if len(req.ID) <= 0 || len(req.Password) <= 0 {
			return fiber.NewError(fiber.StatusBadRequest, "missing username or password field")
		}

		if err = ctrl.db.Where("slug = ?", req.ClientID).First(&client).Error; err != nil {
			return hyperutils.ErrorParser(err)
		} else if client.Secret != req.ClientSecret {
			return fiber.NewError(fiber.StatusForbidden, "invalid client secret")
		}

		user, err = ctrl.auth.AuthUser(req.ID, req.Password)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid username or password")
		}

		exp := time.Now().Add(7 * 24 * time.Hour)
		session = models.UserSession{
			ExpiredAt: &exp,
			Type:      models.UserSessionTypeOauth,
			Scope:     datatypes.NewJSONSlice(strings.Split(req.Scope, " ")),
			IpAddress: c.IP(),
			Location:  "Unknown",
			Available: true,
			ClientID:  &client.ID,
			UserID:    user.ID,
		}

		if err = ctrl.db.Save(&session).Error; err != nil {
			return hyperutils.ErrorParser(err)
		}
	}

	access, accessToken, err := ctrl.auth.NewJwt(session, models.UserClaimsTypeAccess, client.Slug)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	refresh, refreshToken, err := ctrl.auth.NewJwt(session, models.UserClaimsTypeRefresh, client.Slug)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	session.Access = access.ID
	session.Refresh = refresh.ID

	if err := ctrl.db.Save(&session).Error; err != nil {
		return hyperutils.ErrorParser(err)
	} else {
		c.Cookie(&fiber.Cookie{
			Path:     "/",
			Name:     "authorization",
			Value:    accessToken,
			MaxAge:   int(viper.GetDuration("security.sessions_alive_duration").Seconds()),
			SameSite: "None",
			Secure:   true,
		})

		return c.JSON(fiber.Map{
			"id_token":      accessToken,
			"access_token":  accessToken,
			"refresh_token": refreshToken,
			"token_type":    "Bearer",
			"expires_in":    int64(viper.GetDuration("security.sessions_alive_duration").Seconds()),
			"scope":         req.Scope,
			"redirect_uri":  req.Redirect,
		})
	}
}
