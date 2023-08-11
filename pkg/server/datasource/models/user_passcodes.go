package models

import (
	"time"

	"gorm.io/datatypes"
)

const (
	OneTimeVerifyContactCode = iota
	OneTimeDangerousPassCode = iota
)

type OTP struct {
	Model

	Type        int                            `json:"type"`
	Code        string                         `json:"code" gorm:"uniqueIndex"`
	Payload     datatypes.JSONType[OTPPayload] `json:"payload"`
	RefreshedAt *time.Time                     `json:"refreshed_at"`
	ExpiredAt   *time.Time                     `json:"expired_at"`
	UserID      uint                           `json:"user_id"`
}

type OTPPayload struct {
	Target    string `json:"target"`
	IpAddress string `json:"ip_address"`
}
