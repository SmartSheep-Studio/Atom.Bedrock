package services

import (
	models "code.smartsheep.studio/atom/bedrock/pkg/server/datasource/models"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/datatypes"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	db            *gorm.DB
	conf          *viper.Viper
	notifications *NotificationService
}

func NewUserService(db *gorm.DB, conf *viper.Viper, notifications *NotificationService) *UserService {
	return &UserService{db, conf, notifications}
}

func (v *UserService) LookupUser(id string) (models.User, error) {
	var user models.User
	if err := v.db.Where("name = ?", id).Preload("Locks").First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			var contact models.UserContact
			if err := v.db.Where("content = ?", id).First(&contact).Error; err != nil {
				return user, err
			} else if err := v.db.Where("id = ?", contact.UserID).Preload("Locks").First(&user).Error; err != nil {
				return user, err
			}
			return user, nil
		} else {
			return user, err
		}
	} else {
		return user, nil
	}
}

func (v *UserService) NewUser(item *models.User) error {
	encrypted, _ := bcrypt.GenerateFromPassword([]byte(item.Password), bcrypt.DefaultCost)
	item.Password = string(encrypted)

	if item.Permissions.Data() == nil || len(item.Permissions.Data()) == 0 {
		item.Permissions = datatypes.NewJSONType(v.conf.GetStringMap("security.preset_permissions.default"))
	}

	if item.VerifiedAt == nil {
		if err := v.notifications.SendNotification(&models.Notification{
			Title:       "Account verification is required.",
			Description: "Don't forgot verify your account!",
			Content:     "Your account isn't verified now. Before you verify, some features are unavailable. Now go to account center and verify your account now!",
			Level:       models.NotificationLevelWarning,
			SenderType:  models.NotificationSenderTypeSystem,
			SenderID:    nil,
		}); err != nil {
			return err
		}
	}

	if err := v.notifications.SendNotification(&models.Notification{
		Title:       fmt.Sprintf("Welcome to %s", v.conf.GetString("general.name")),
		Description: fmt.Sprintf("Thanks for you choosing %s.", v.conf.GetString("general.name")),
		Content:     fmt.Sprintf("Thanks for you registration of %s. Now go to explore the whole platform!", v.conf.GetString("general.name")),
		Level:       models.NotificationLevelInfo,
		SenderType:  models.NotificationSenderTypeSystem,
		SenderID:    nil,
	}); err != nil {
		return err
	}

	return v.db.Save(&item).Error
}
