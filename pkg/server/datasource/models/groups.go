package models

import "gorm.io/datatypes"

type UserGroup struct {
	Model

	Slug        string                             `json:"slug" gorm:"uniqueIndex"`
	Name        string                             `json:"name"`
	Description string                             `json:"description"`
	Users       []User                             `json:"users" gorm:"many2many:user_joined_groups"`
	Permissions datatypes.JSONType[map[string]any] `json:"permissions"`
}
