package models

import (
	"code.smartsheep.studio/atom/bedrock/pkg/kit/common"
	"fmt"
	"gorm.io/gorm"
	"time"

	"gorm.io/datatypes"
)

type User struct {
	Model

	AvatarUrl     string                             `json:"avatar_url"`
	BannerUrl     string                             `json:"banner_url"`
	Name          string                             `json:"name" gorm:"uniqueIndex"`
	Nickname      string                             `json:"nickname"`
	Password      string                             `json:"-"`
	Description   string                             `json:"description"`
	Contacts      []UserContact                      `json:"contacts"`
	Notifications []Notification                     `json:"notifications" gorm:"foreignKey:RecipientID"`
	OauthClients  []OauthClient                      `json:"oauth_clients"`
	Sessions      []UserSession                      `json:"sessions"`
	PassCodes     []OTP                              `json:"passcodes"`
	Assets        []StorageFile                      `json:"user_assets"`
	Friends       []*User                            `json:"friends" gorm:"many2many:user_friends"`
	Groups        []UserGroup                        `json:"groups" gorm:"many2many:user_joined_groups"`
	Permissions   datatypes.JSONType[map[string]any] `json:"permissions"`
	VerifiedAt    *time.Time                         `json:"verified_at"`
	LockedAt      *time.Time                         `json:"locked_at"`
}

//goland:noinspection GoMixedReceiverTypes
func (v User) GetPermissions() (map[string]any, error) {
	perms := map[string]any{}
	fmt.Println(v.Groups)
	for _, group := range v.Groups {
		// Merge into permissions map
		for k, val := range group.Permissions.Data() {
			perms[k] = val
		}
	}
	// User self's permissions will override group permissions
	for k, val := range v.Permissions.Data() {
		perms[k] = val
	}
	// Return
	return perms, nil
}

//goland:noinspection GoMixedReceiverTypes
func (v User) HasPermissions(requires ...string) error {
	if perms, err := v.GetPermissions(); err != nil {
		return err
	} else {
		return common.MatchTree(perms, requires...)
	}
}

//goland:noinspection GoMixedReceiverTypes
func (v *User) BeforeCreate(tx *gorm.DB) error {
	if len(v.Permissions.Data()) == 0 {
		v.Permissions = datatypes.NewJSONType(map[string]any{})
	}

	return nil
}

const (
	UserContactTypeEmail = "email"
	UserContactTypePhone = "phone"
)

type UserContact struct {
	Model

	Type        string     `json:"type"`
	Name        string     `json:"name"`
	Content     string     `json:"content" gorm:"uniqueIndex"`
	Description string     `json:"description"`
	VerifiedAt  *time.Time `json:"verified_at"`
	IsPrimary   bool       `json:"is_primary"`
	IsSecondary bool       `json:"is_secondary"`
	UserID      uint       `json:"user_id"`
}

type UserGroup struct {
	Model

	Slug        string                             `json:"slug" gorm:"uniqueIndex"`
	Name        string                             `json:"name"`
	Description string                             `json:"description"`
	Users       []User                             `json:"users" gorm:"many2many:user_joined_groups"`
	Permissions datatypes.JSONType[map[string]any] `json:"permissions"`
}
