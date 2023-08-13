package models

import "time"

type Lock struct {
	Model

	Reason    string     `json:"reason"`
	ExpiredAt *time.Time `json:"expired_at"`
	UserID    *uint      `json:"user_id"`
}
