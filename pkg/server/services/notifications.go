package services

import (
	"code.smartsheep.studio/atom/bedrock/pkg/server/datasource/models"
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type NotificationService struct {
	db     *gorm.DB
	mailer *MailerService
}

func NewNotificationService(db *gorm.DB, mailer *MailerService) *NotificationService {
	return &NotificationService{db, mailer}
}

func (v *NotificationService) SendNotification(item *models.Notification) error {
	var recipient models.User
	if err := v.db.Where("id = ?", item.RecipientID).Preload("Contacts").First(&recipient).Error; err != nil {
		return err
	}

	contact, ok := recipient.GetPrimaryContact()
	if ok && item.Level != "tips" {
		return v.mailer.SendMail(
			contact.Content,
			fmt.Sprintf("⌈%s · Notification⌋ %s", viper.GetString("general.name"), item.Title),
			fmt.Sprintf(`Hello, %s!
You have received a new notification.

%s
%s

%s

Notification ID: #%d
Created At: %s`, recipient.Nickname, item.Title, item.Description, item.Content, item.ID, item.CreatedAt),
		)
	}

	return v.db.Save(&item).Error
}
