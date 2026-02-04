package models

import (
	"time"

	"gorm.io/gorm"
)

type Trainer struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:255;not null" json:"name"`
	Bio       string         `gorm:"type:text" json:"bio"`
	Specialty string         `gorm:"size:255" json:"specialty"`
	PhotoURL  string         `gorm:"size:500" json:"photo_url"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Trainer) TableName() string {
	return "trainers"
}

type TrainerSchedule struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	TrainerID uint           `gorm:"index;not null" json:"trainer_id"`
	DayOfWeek int            `gorm:"not null" json:"day_of_week"`       // 0=Sunday ... 6=Saturday
	StartTime string         `gorm:"size:5;not null" json:"start_time"` // HH:MM
	EndTime   string         `gorm:"size:5;not null" json:"end_time"`   // HH:MM
	Location  string         `gorm:"size:255" json:"location"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (TrainerSchedule) TableName() string {
	return "trainer_schedules"
}
