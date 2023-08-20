package controllers

import (
	models "code.smartsheep.studio/atom/bedrock/pkg/server/datasource/models"
	hyperutils "code.smartsheep.studio/atom/bedrock/pkg/server/hypertext/hyperutils"
	"code.smartsheep.studio/atom/bedrock/pkg/server/hypertext/middlewares"
	"code.smartsheep.studio/atom/bedrock/pkg/server/services"
	"net/http"
	"os"
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
	router.Get(
		"/api/assets/:id/meta",
		ctrl.readMeta,
	)

	router.Post(
		"/cgi/assets",
		ctrl.gatekeeper.Fn(false, hyperutils.GenScope("create:assets"), hyperutils.GenPerms("assets.upload")),
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

	var item models.StorageFile
	if err := tx.First(&item).Error; err != nil {
		return hyperutils.ErrorParser(err)
	} else {
		return c.Download(filepath.Join(viper.GetString("paths.user_contents"), item.StorageID), item.Name)
	}
}

func (ctrl *StorageController) readMeta(c *fiber.Ctx) error {
	probe := c.Params("id")

	var tx *gorm.DB
	if _, err := strconv.Atoi(probe); err == nil {
		tx = ctrl.db.Where("id = ?", probe)
	} else {
		tx = ctrl.db.Where("storage_id = ?", probe)
	}

	var item models.StorageFile
	if err := tx.First(&item).Error; err != nil {
		return hyperutils.ErrorParser(err)
	}

	if data, err := os.ReadFile(filepath.Join(viper.GetString("paths.user_contents"), item.StorageID)); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	} else {
		return c.JSON(fiber.Map{
			"record":   item,
			"mimetype": http.DetectContentType(data),
			"size":     item.Size,
		})
	}
}

func (ctrl *StorageController) upload(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if ok := c.Locals("principal-ok").(bool); ok {
		user := c.Locals("principal").(models.User)
		if item, err := ctrl.warehouse.SaveFile2User(c, file, user, models.StorageFileCustomType); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		} else {
			return c.JSON(item)
		}
	} else {
		if item, err := ctrl.warehouse.SaveFile(c, file, models.StorageFileCustomType); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		} else {
			return c.JSON(item)
		}
	}
}
