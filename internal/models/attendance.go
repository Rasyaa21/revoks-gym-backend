package models

import (
	"time"

	"gorm.io/gorm"
)

type AttendanceLog struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	UserID     uint           `gorm:"index;not null" json:"user_id"`
	Direction  string         `gorm:"size:10;not null" json:"direction"`           // in|out
	Source     string         `gorm:"size:20;not null;default:'qr'" json:"source"` // qr|manual
	OccurredAt time.Time      `gorm:"not null" json:"occurred_at"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (AttendanceLog) TableName() string {
	return "attendance_logs"
}
