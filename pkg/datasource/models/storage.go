package models

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"gorm.io/gorm"
)

const (
	StorageFileCustomType = iota
	StorageFileBannerType
	StorageFileAvatarType
	StorageFileProfileType
)

type StorageFile struct {
	Model

	Name      string `json:"name"`
	Size      int64  `json:"size"`
	Type      int    `json:"type"`
	StorageID string `json:"storage_id"`
	UserID    *uint  `json:"user_id"`
}

func (v *StorageFile) AfterDelete(tx *gorm.DB) (err error) {
	// Remove file in filesystem
	os.Remove(filepath.Join(viper.GetString("paths.user_contents"), v.StorageID))
	return nil
}

func (v StorageFile) GetURL() string {
	return fmt.Sprintf("/api/assets/%s", v.StorageID)
}
