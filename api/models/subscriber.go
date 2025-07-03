package models

import "time"

// Subscriber represents a single subscriber record.
type Subscriber struct {
	ID               uint       `gorm:"primaryKey" json:"id"`
	Email            string     `gorm:"type:varchar(255);not null;unique" json:"email"`
	Name             string     `gorm:"type:varchar(255);not null" json:"name"`
	Newsletter       *bool      `gorm:"default:false" json:"newsletter"`
	EmailValidatedAt *time.Time `gorm:"column:email_validated_at" json:"email_validated_at,omitempty"`
	CreatedAt        time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Subscriber) TableName() string { return "api.subscriber" }
