package models

import (
	"gorm.io/datatypes"
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

type NotificationLink struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
}
