package models

import (
	"time"

	"gorm.io/gorm"
)

type Membership struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"index;not null" json:"user_id"`
	Status    string         `gorm:"size:20;not null;default:'active'" json:"status"` // active|expired
	Plan      string         `gorm:"size:50" json:"plan"`
	StartsAt  time.Time      `gorm:"not null" json:"starts_at"`
	EndsAt    time.Time      `gorm:"not null" json:"ends_at"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Membership) TableName() string {
	return "memberships"
}
