package models

import "time"

type LoginAttempt struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    *uint     `gorm:"index"              json:"user_id,omitempty"`
	Email     string    `gorm:"size:255;index"     json:"email"`
	Success   bool      `gorm:"index;not null"    json:"success"`
	Reason    string    `gorm:"size:64"           json:"reason,omitempty"`
	IP        *string   `gorm:"size:45"         json:"ip"`
	UserAgent string    `gorm:"size:255"        json:"user_agent"`
	CreatedAt time.Time `json:"created_at"`
}

func (LoginAttempt) TableName() string { return "api.login_attempt" }
