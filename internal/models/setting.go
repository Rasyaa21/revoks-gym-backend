package models

import (
	"time"

	"gorm.io/gorm"
)

type UserSetting struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	UserID       uint           `gorm:"uniqueIndex;not null" json:"user_id"`
	PushEnabled  bool           `gorm:"not null;default:true" json:"push_enabled"`
	EmailEnabled bool           `gorm:"not null;default:true" json:"email_enabled"`
	Language     string         `gorm:"size:10;not null;default:'id'" json:"language"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (UserSetting) TableName() string {
	return "user_settings"
}
