package controllers

import (
	"code.smartsheep.studio/atom/bedrock/pkg/services"
	"path/filepath"

	"code.smartsheep.studio/atom/bedrock/pkg/hypertext/hyperutils"
	"github.com/gofiber/fiber/v2"

	"code.smartsheep.studio/atom/bedrock/pkg/datasource/models"
	"code.smartsheep.studio/atom/bedrock/pkg/hypertext/middlewares"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type StorageController struct {
	db         *gorm.DB
	warehouse  *services.StorageService
	gatekeeper *middlewares.AuthMiddleware
}

func NewStorageController(db *gorm.DB, warehouse *services.StorageService, gatekeeper *middlewares.AuthMiddleware) *StorageController {
	ctrl := &StorageController{db, warehouse, gatekeeper}
	return ctrl
}

func (ctrl *StorageController) Map(router *fiber.App) {
	router.Get(
		"/api/assets",
		ctrl.read,
	)
	router.Post(
		"/api/assets",
		ctrl.gatekeeper.Fn(true, hyperutils.GenScope("create:assets"), hyperutils.GenPerms("assets.upload")),
		ctrl.upload,
	)
}

func (ctrl *StorageController) read(c *fiber.Ctx) error {
	id := c.Query("id")
	storage := c.Query("storage")

	if len(id) == 0 {
		id = "0"
	}
	if len(storage) == 0 {
		storage = "?"
	}

	var f models.StorageFile
	if err := ctrl.db.Where("id = ? OR storage_id = ?", id, storage).First(&f).Error; err != nil {
		return hyperutils.ErrorParser(err)
	} else {
		return c.Download(filepath.Join(viper.GetString("paths.user_contents"), f.StorageID), f.Name)
	}
}

func (ctrl *StorageController) upload(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.User)

	file, err := c.FormFile("file")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if f, err := ctrl.warehouse.SaveFile2User(c, file, user, models.StorageFileCustomType); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	} else {
		return c.JSON(f)
	}
}
