package models

import (
	"time"

	"gorm.io/gorm"
)

type Template struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:255;not null" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Template) TableName() string {
	return "templates"
}

type UserTemplateFollow struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	UserID     uint           `gorm:"index;not null" json:"user_id"`
	TemplateID uint           `gorm:"index;not null" json:"template_id"`
	FollowedAt time.Time      `gorm:"not null" json:"followed_at"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (UserTemplateFollow) TableName() string {
	return "user_template_follows"
}
