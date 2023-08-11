package controllers

import (
	models "code.smartsheep.studio/atom/bedrock/pkg/server/datasource/models"
	hyperutils "code.smartsheep.studio/atom/bedrock/pkg/server/hypertext/hyperutils"
	"code.smartsheep.studio/atom/bedrock/pkg/server/hypertext/middlewares"
	"code.smartsheep.studio/atom/bedrock/pkg/server/services"
	"github.com/gofiber/fiber/v2"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type OauthController struct {
	db         *gorm.DB
	auth       *services.AuthService
	gatekeeper *middlewares.AuthMiddleware
}

func NewOauthController(db *gorm.DB, auth *services.AuthService, gatekeeper *middlewares.AuthMiddleware) *OauthController {
	return &OauthController{db, auth, gatekeeper}
}

func (ctrl *OauthController) Map(router *fiber.App) {
	router.Get(
		"/api/users/oauth",
		ctrl.gatekeeper.Fn(true, hyperutils.GenScope("read:oauth"), hyperutils.GenPerms("users.oauth.read")),
		ctrl.list,
	)
	router.Get(
		"/api/users/oauth/:oauth",
		ctrl.gatekeeper.Fn(true, hyperutils.GenScope("read:oauth"), hyperutils.GenPerms("users.oauth.read")),
		ctrl.get,
	)
	router.Post(
		"/api/users/oauth",
		ctrl.gatekeeper.Fn(true, hyperutils.GenScope("create:oauth"), hyperutils.GenPerms("users.oauth.create")),
		ctrl.create,
	)
	router.Put(
		"/api/users/oauth/:oauth",
		ctrl.gatekeeper.Fn(true, hyperutils.GenScope("update:oauth"), hyperutils.GenPerms("users.oauth.update")),
		ctrl.update,
	)
	router.Delete(
		"/api/users/oauth/:oauth",
		ctrl.gatekeeper.Fn(true, hyperutils.GenScope("delete:oauth"), hyperutils.GenPerms("users.oauth.delete")),
		ctrl.delete,
	)
}

func (ctrl *OauthController) list(c *fiber.Ctx) error {
	u := c.Locals("principal").(models.User)

	var clients []models.OauthClient
	if err := ctrl.db.Where("user_id = ?", u.ID).Find(&clients).Error; err != nil {
		return hyperutils.ErrorParser(err)
	} else {
		return c.JSON(clients)
	}
}

func (ctrl *OauthController) get(c *fiber.Ctx) error {
	u := c.Locals("principal").(models.User)

	var client models.OauthClient
	if err := ctrl.db.Where("user_id = ? AND slug = ?", u.ID, c.Params("oauth")).First(&client).Error; err != nil {
		return hyperutils.ErrorParser(err)
	} else {
		return c.JSON(client)
	}
}

func (ctrl *OauthController) create(c *fiber.Ctx) error {
	u := c.Locals("principal").(models.User)

	var req struct {
		Slug        string   `json:"slug" validate:"required"`
		Name        string   `json:"name" validate:"required"`
		Description string   `json:"description"`
		Secret      string   `json:"secret" validate:"required"`
		Urls        []string `json:"urls"`
		Callbacks   []string `json:"callbacks"`
	}

	if err := hyperutils.BodyParser(c, &req); err != nil {
		return err
	}

	client := models.OauthClient{
		Slug:         req.Slug,
		Name:         req.Name,
		Description:  req.Description,
		Secret:       req.Secret,
		Urls:         datatypes.NewJSONSlice(req.Urls),
		Callbacks:    datatypes.NewJSONSlice(req.Callbacks),
		IsDeveloping: true,
		UserID:       &u.ID,
	}

	if err := ctrl.db.Save(&client).Error; err != nil {
		return hyperutils.ErrorParser(err)
	} else {
		return c.JSON(client)
	}
}

func (ctrl *OauthController) update(c *fiber.Ctx) error {
	u := c.Locals("principal").(models.User)

	var req struct {
		Slug        string   `json:"slug" validate:"required"`
		Name        string   `json:"name" validate:"required"`
		Description string   `json:"description"`
		Secret      string   `json:"secret" validate:"required"`
		Urls        []string `json:"urls"`
		Callbacks   []string `json:"callbacks"`
	}

	if err := hyperutils.BodyParser(c, &req); err != nil {
		return err
	}

	var client models.OauthClient
	if err := ctrl.db.Where("user_id = ? AND slug = ?", u.ID, c.Params("oauth")).First(&client).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}

	client.Slug = req.Slug
	client.Name = req.Name
	client.Description = req.Description
	client.Secret = req.Secret
	client.Urls = datatypes.NewJSONSlice(req.Urls)
	client.Callbacks = datatypes.NewJSONSlice(req.Callbacks)

	if err := ctrl.db.Save(&client).Error; err != nil {
		return hyperutils.ErrorParser(err)
	} else {
		return c.JSON(client)
	}
}

func (ctrl *OauthController) delete(c *fiber.Ctx) error {
	u := c.Locals("principal").(models.User)

	var client models.OauthClient
	if err := ctrl.db.Where("user_id = ? AND slug = ?", u.ID, c.Params("oauth")).First(&client).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}

	if err := ctrl.db.Delete(&client).Error; err != nil {
		return hyperutils.ErrorParser(err)
	} else {
		return c.SendStatus(fiber.StatusNoContent)
	}
}
