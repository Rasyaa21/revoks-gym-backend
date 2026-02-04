package models

import (
	"time"

	"gorm.io/gorm"
)

type Target struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"index;not null" json:"user_id"`
	Period    string         `gorm:"size:20;not null" json:"period"` // weekly|monthly
	Title     string         `gorm:"size:255;not null" json:"title"`
	GoalValue int            `gorm:"not null" json:"goal_value"`
	StartDate time.Time      `gorm:"not null" json:"start_date"`
	EndDate   time.Time      `gorm:"not null" json:"end_date"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Target) TableName() string {
	return "targets"
}

type TargetProgress struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	TargetID   uint           `gorm:"index;not null" json:"target_id"`
	Value      int            `gorm:"not null" json:"value"`
	RecordedAt time.Time      `gorm:"not null" json:"recorded_at"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (TargetProgress) TableName() string {
	return "target_progress"
}
