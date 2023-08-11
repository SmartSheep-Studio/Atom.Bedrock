package controllers

import (
	models2 "code.smartsheep.studio/atom/bedrock/pkg/server/datasource/models"
	hyperutils2 "code.smartsheep.studio/atom/bedrock/pkg/server/hypertext/hyperutils"
	"code.smartsheep.studio/atom/bedrock/pkg/server/hypertext/middlewares"
	"code.smartsheep.studio/atom/bedrock/pkg/server/services"
	"path/filepath"
	"strconv"

	"github.com/gofiber/fiber/v2"

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
		"/api/assets/:id",
		ctrl.read,
	)
	router.Post(
		"/api/assets",
		ctrl.gatekeeper.Fn(true, hyperutils2.GenScope("create:assets"), hyperutils2.GenPerms("assets.upload")),
		ctrl.upload,
	)
}

func (ctrl *StorageController) read(c *fiber.Ctx) error {
	probe := c.Params("id")

	var tx *gorm.DB
	if _, err := strconv.Atoi(probe); err == nil {
		tx = ctrl.db.Where("id = ?", probe)
	} else {
		tx = ctrl.db.Where("storage_id = ?", probe)
	}

	var f models2.StorageFile
	if err := tx.First(&f).Error; err != nil {
		return hyperutils2.ErrorParser(err)
	} else {
		return c.Download(filepath.Join(viper.GetString("paths.user_contents"), f.StorageID), f.Name)
	}
}

func (ctrl *StorageController) upload(c *fiber.Ctx) error {
	user := c.Locals("principal").(models2.User)

	file, err := c.FormFile("file")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if f, err := ctrl.warehouse.SaveFile2User(c, file, user, models2.StorageFileCustomType); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	} else {
		return c.JSON(f)
	}
}
