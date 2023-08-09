package controllers

import (
	"code.smartsheep.studio/atom/bedrock/pkg/hypertext/hyperutils"
	"code.smartsheep.studio/atom/bedrock/pkg/services"
	"fmt"
	"github.com/gofiber/fiber/v2"

	"code.smartsheep.studio/atom/bedrock/pkg/datasource/models"
	"code.smartsheep.studio/atom/bedrock/pkg/hypertext/middlewares"
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
		ctrl.gatekeeper.Fn(true, hyperutils.GenScope("read:id.tokens"), hyperutils.GenPerms()),
		ctrl.list,
	)
	router.Post(
		"/api/users/tokens",
		ctrl.gatekeeper.Fn(true, hyperutils.GenScope("create:id.tokens"), hyperutils.GenPerms()),
		ctrl.create,
	)
}

func (ctrl *ApiTokenController) list(c *fiber.Ctx) error {
	u := c.Locals("principal").(models.User)

	var tokens []models.UserSession
	if err := ctrl.db.Where("user_id = ? AND type = ?", u.ID, models.UserSessionTypeToken).Find(&tokens).Error; err != nil {
		return hyperutils.ErrorParser(err)
	} else {
		return c.JSON(tokens)
	}
}

func (ctrl *ApiTokenController) create(c *fiber.Ctx) error {
	u := c.Locals("principal").(models.User)

	var req struct {
		Description string   `json:"description" validate:"required"`
		Scope       []string `json:"scope" validate:"required"`
	}

	if err := hyperutils.BodyParser(c, &req); err != nil {
		return err
	}

	session := models.UserSession{
		IpAddress:   c.IP(),
		Type:        models.UserSessionTypeToken,
		UserID:      u.ID,
		Scope:       datatypes.NewJSONSlice(req.Scope),
		Description: req.Description,
		ExpiredAt:   nil,
		Available:   true,
	}

	if err := ctrl.db.Save(&session).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}

	_, token, err := ctrl.auth.NewJwt(session, models.UserClaimsTypeAccess)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to encode jwt: %s", err.Error()))
	} else {
		return c.JSON(fiber.Map{
			"session": session,
			"token":   token,
		})
	}
}
