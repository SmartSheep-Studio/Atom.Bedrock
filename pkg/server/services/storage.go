package services

import (
	models "code.smartsheep.studio/atom/bedrock/pkg/server/datasource/models"
	"github.com/gofiber/fiber/v2"
	"mime/multipart"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type StorageService struct {
	db *gorm.DB
}

func NewStorageService(db *gorm.DB) *StorageService {
	return &StorageService{db: db}
}

func (v *StorageService) SaveFile(c *fiber.Ctx, file *multipart.FileHeader, flag int) (models.StorageFile, error) {
	f := models.StorageFile{
		Name:      file.Filename,
		Size:      file.Size,
		Type:      flag,
		StorageID: uuid.NewString(),
		UserID:    nil,
	}

	if err := v.db.Save(&f).Error; err != nil {
		return f, err
	} else {
		file.Filename = f.StorageID
		return f, c.SaveFile(file, filepath.Join(viper.GetString("paths.user_contents"), file.Filename))
	}
}

func (v *StorageService) SaveFile2User(c *fiber.Ctx, file *multipart.FileHeader, user models.User, flag int) (models.StorageFile, error) {
	f := models.StorageFile{
		Name:      file.Filename,
		Size:      file.Size,
		Type:      flag,
		StorageID: uuid.NewString(),
		UserID:    &user.ID,
	}

	if err := v.db.Save(&f).Error; err != nil {
		return f, err
	} else {
		file.Filename = f.StorageID
		return f, c.SaveFile(file, filepath.Join(viper.GetString("paths.user_contents"), file.Filename))
	}
}
