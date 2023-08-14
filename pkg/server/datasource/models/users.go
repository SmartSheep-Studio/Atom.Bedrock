package models

import (
	"fmt"
	"time"

	"code.smartsheep.studio/atom/bedrock/pkg/kit/common"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"

	"gorm.io/datatypes"
)

type User struct {
	Model

	AvatarUrl          string                             `json:"avatar_url"`
	BannerUrl          string                             `json:"banner_url"`
	Name               string                             `json:"name" gorm:"uniqueIndex"`
	Nickname           string                             `json:"nickname"`
	Password           string                             `json:"-"`
	Description        string                             `json:"description"`
	Contacts           []UserContact                      `json:"contacts"`
	Notifications      []Notification                     `json:"notifications" gorm:"foreignKey:RecipientID"`
	NotificationsCount int64                              `json:"notifications_count" gorm:"-"`
	OauthClients       []OauthClient                      `json:"oauth_clients"`
	Sessions           []UserSession                      `json:"sessions"`
	PassCodes          []OTP                              `json:"passcodes"`
	Assets             []StorageFile                      `json:"user_assets"`
	Friends            []*User                            `json:"friends" gorm:"many2many:user_friends"`
	Groups             []UserGroup                        `json:"groups" gorm:"many2many:user_joined_groups"`
	Locks              []Lock                             `json:"locks"`
	Permissions        datatypes.JSONType[map[string]any] `json:"permissions"`
	VerifiedAt         *time.Time                         `json:"verified_at"`
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

const (
	UserSessionTypeOauth = iota
	UserSessionTypeAuth
	UserSessionTypeToken
)

type UserSession struct {
	Model

	IpAddress   string                      `json:"ip"`
	Location    string                      `json:"location"`
	Available   bool                        `json:"available"`
	Type        int                         `json:"type"`
	Code        string                      `json:"code" gorm:"type:varchar(512)"`
	Access      string                      `json:"access" gorm:"type:varchar(512)"`
	Refresh     string                      `json:"refresh" gorm:"type:varchar(512)"`
	Description string                      `json:"description"`
	Scope       datatypes.JSONSlice[string] `json:"scope"`
	ExpiredAt   *time.Time                  `json:"expired_at"`
	ClientID    *uint                       `json:"client_id"`
	UserID      uint                        `json:"user_id"`
}

// HasScope use non ptr receiver because it usually used in non ptr model
//
//goland:noinspection GoMixedReceiverTypes
func (u UserSession) HasScope(requires ...string) error {
	return common.MatchList(u.Scope, requires...)
}

// BeforeCreate is a gorm hook
//
//goland:noinspection GoMixedReceiverTypes
func (u *UserSession) BeforeCreate(tx *gorm.DB) (err error) {
	u.Location = "Unknown"

	return nil
}

const (
	UserClaimsTypeAccess  = "access_token"
	UserClaimsTypeRefresh = "refresh_token"
)

type UserClaims struct {
	jwt.RegisteredClaims

	Type            string `json:"typ"`
	SessionID       uint   `json:"session_id"`
	ClientID        *uint  `json:"client_id"`
	PersonalTokenID *uint  `json:"token_id"`
}
