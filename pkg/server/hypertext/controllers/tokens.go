package controllers

import (
	models2 "code.smartsheep.studio/atom/bedrock/pkg/server/datasource/models"
	hyperutils2 "code.smartsheep.studio/atom/bedrock/pkg/server/hypertext/hyperutils"
	"code.smartsheep.studio/atom/bedrock/pkg/server/hypertext/middlewares"
	"code.smartsheep.studio/atom/bedrock/pkg/server/services"
	"fmt"
	"github.com/gofiber/fiber/v2"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ApiTokenController struct {
	db         *gorm.DB
	auth       *services.AuthService
	gatekeeper *middlewares.AuthMiddleware
}

func NewApiTokenController(db *gorm.DB, auth *services.AuthService, gatekeeper *middlewares.AuthMiddleware) *ApiTokenController {
	ctrl := &ApiTokenController{db, auth, gatekeeper}
	return ctrl
}

func (ctrl *ApiTokenController) Map(router *fiber.App) {
	router.Get(
		"/api/users/tokens",
		ctrl.gatekeeper.Fn(true, hyperutils2.GenScope("read:id.tokens"), hyperutils2.GenPerms()),
		ctrl.list,
	)
	router.Post(
		"/api/users/tokens",
		ctrl.gatekeeper.Fn(true, hyperutils2.GenScope("create:id.tokens"), hyperutils2.GenPerms()),
		ctrl.create,
	)
}

func (ctrl *ApiTokenController) list(c *fiber.Ctx) error {
	u := c.Locals("principal").(models2.User)

	var tokens []models2.UserSession
	if err := ctrl.db.Where("user_id = ? AND type = ?", u.ID, models2.UserSessionTypeToken).Find(&tokens).Error; err != nil {
		return hyperutils2.ErrorParser(err)
	} else {
		return c.JSON(tokens)
	}
}

func (ctrl *ApiTokenController) create(c *fiber.Ctx) error {
	u := c.Locals("principal").(models2.User)

	var req struct {
		Description string   `json:"description" validate:"required"`
		Scope       []string `json:"scope" validate:"required"`
	}

	if err := hyperutils2.BodyParser(c, &req); err != nil {
		return err
	}

	session := models2.UserSession{
		IpAddress:   c.IP(),
		Type:        models2.UserSessionTypeToken,
		UserID:      u.ID,
		Scope:       datatypes.NewJSONSlice(req.Scope),
		Description: req.Description,
		ExpiredAt:   nil,
		Available:   true,
	}

	if err := ctrl.db.Save(&session).Error; err != nil {
		return hyperutils2.ErrorParser(err)
	}

	_, token, err := ctrl.auth.NewJwt(session, models2.UserClaimsTypeAccess)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to encode jwt: %s", err.Error()))
	} else {
		return c.JSON(fiber.Map{
			"session": session,
			"token":   token,
		})
	}
}
