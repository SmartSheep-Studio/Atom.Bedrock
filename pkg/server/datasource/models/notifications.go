package models

import (
	"code.smartsheep.studio/atom/bedrock/pkg/server/services"
	"fmt"
	"github.com/samber/lo"
	"github.com/spf13/viper"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

const (
	// NotificationLevelTips Won't send alert
	NotificationLevelTips = "tips"
	// NotificationLevelInfo Will send alert
	NotificationLevelInfo = "info"
	// NotificationLevelWarning Will send alert
	NotificationLevelWarning = "warning"
	// NotificationLevelAlert Will send alert
	NotificationLevelAlert = "alert"
)

const (
	NotificationSenderTypeSystem = "system"
	NotificationSenderTypeUser   = "user"
)

type Notification struct {
	Model

	Title       string                                `json:"title"`
	Description string                                `json:"description"`
	Content     string                                `json:"content"`
	Level       string                                `json:"level"`
	Links       datatypes.JSONSlice[NotificationLink] `json:"links"`
	ReadAt      *time.Time                            `json:"read_at"`
	RecipientID uint                                  `json:"recipient_id"`
	SenderType  string                                `json:"sender_type"`
	SenderID    *uint                                 `json:"sender_id"`
}

func (v *Notification) AfterCreate(tx *gorm.DB) error {
	var recipient User
	if err := tx.Where("id = ?", v.RecipientID).Preload("Contacts").First(&recipient).Error; err != nil {
		return err
	}

	contact, ok := recipient.GetPrimaryContact()
	if ok && !lo.Contains([]string{"tips"}, v.Level) {
		services.Mailer.SendMail(
			contact.Content,
			fmt.Sprintf("⌈%s · Notification⌋ %s", viper.GetString("general.name"), v.Title),
			fmt.Sprintf(`Hello, %s!
You have received a new notification.

%s
%s

%s

Notification ID: #%d
Created At: %s`, recipient.Nickname, v.Title, v.Description, v.Content, v.ID, v.CreatedAt),
		)
	}

	return nil
}

type NotificationLink struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
}
