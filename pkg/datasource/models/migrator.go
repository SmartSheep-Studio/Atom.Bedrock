package models

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func RunMigration(db *gorm.DB) {
	// Creating tables
	if err := db.AutoMigrate(
		&Notification{},
		&UserGroup{},
		&OauthClient{},
		&User{},
		&UserContact{},
		&UserSession{},
		&OTP{},
		&StorageFile{},
	); err != nil {
		log.Fatal().Err(err).Msg("Error when migrating database")
	}

	// Seeding data
	if err := db.First(&UserGroup{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		db.Create([]*UserGroup{
			{
				Slug:        "administrators",
				Name:        "Administrators",
				Description: "Have all permissions users.",
				Permissions: datatypes.NewJSONType(viper.GetStringMap("security.groups.administrators")),
			},
			{
				Slug:        "verified_users",
				Name:        "Verified Users",
				Description: "Have least one verified contact users.",
				Permissions: datatypes.NewJSONType(viper.GetStringMap("security.groups.verified_users")),
			},
		})
	}
	if err := db.First(&User{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		var group UserGroup
		if err := db.Where("slug = ?", "administrators").First(&group).Error; err != nil {
			log.Fatal().Err(err).Msg("Couldn't create default user administrator")
		}
		password := strings.ReplaceAll(uuid.NewString(), "-", "")
		encrypted, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if db.Create([]*User{
			{
				Name:     "administrator",
				Nickname: "Administrator",
				Contacts: []UserContact{
					{
						Name:        "Administrator's Primary Contact",
						Description: "Primary Email",
						Content:     "administrator@example.com",
						Type:        UserContactTypeEmail,
						IsPrimary:   true,
					},
				},
				Password: string(encrypted),
				Groups: []UserGroup{
					group,
				},
			},
		}).Error == nil {
			log.Info().Msgf("Successfully created default user `administrator` with password `%s`", password)
		}
	}
}
