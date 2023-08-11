package models

import (
	"code.smartsheep.studio/atom/bedrock/pkg/kit/common"
	"github.com/golang-jwt/jwt/v5"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

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
