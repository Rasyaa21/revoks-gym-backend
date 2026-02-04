package models

import (
	"time"

	"gorm.io/gorm"
)

type WorkoutProgress struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	UserID     uint           `gorm:"index;not null" json:"user_id"`
	Title      string         `gorm:"size:255;not null" json:"title"`
	Notes      string         `gorm:"type:text" json:"notes"`
	RecordedAt time.Time      `gorm:"not null" json:"recorded_at"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (WorkoutProgress) TableName() string {
	return "workout_progress"
}
