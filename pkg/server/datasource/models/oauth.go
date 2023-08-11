package models

import "gorm.io/datatypes"

type OauthClient struct {
	Model

	Slug         string                      `json:"slug" gorm:"uniqueIndex"`
	Name         string                      `json:"name"`
	Description  string                      `json:"description"`
	Secret       string                      `json:"secret" gorm:"type:varchar(512)"`
	Urls         datatypes.JSONSlice[string] `json:"urls"`
	Callbacks    datatypes.JSONSlice[string] `json:"callbacks"`
	Sessions     []UserSession               `json:"sessions" gorm:"foreignKey:ClientID"`
	IsDanger     bool                        `json:"is_danger"`
	IsOfficial   bool                        `json:"is_official"`
	IsVerified   bool                        `json:"is_verified"`
	IsDeveloping bool                        `json:"is_developing"`
	UserID       *uint                       `json:"user_id"`
}
