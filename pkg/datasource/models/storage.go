package models

import (
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

func (u *StorageFile) AfterDelete(tx *gorm.DB) (err error) {
	// Remove file in filesystem
	os.Remove(filepath.Join(viper.GetString("paths.user_contents"), u.StorageID))
	return nil
}
