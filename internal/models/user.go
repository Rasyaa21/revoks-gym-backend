package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	UserID       uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"user_id"`
	Role         UserRole   `gorm:"type:varchar(20);not null" json:"role"`
	Email        *string    `gorm:"type:varchar(255);uniqueIndex" json:"email"`
	Phone        *string    `gorm:"type:varchar(50);uniqueIndex" json:"phone"`
	PasswordHash string     `gorm:"type:varchar(255);not null" json:"-"`
	FullName     string     `gorm:"type:varchar(255);not null" json:"full_name"`
	PhotoURL     *string    `gorm:"type:varchar(500)" json:"photo_url"`
	Status       UserStatus `gorm:"type:varchar(20);not null;default:'active'" json:"status"`
	CreatedAt    time.Time  `gorm:"type:timestamptz;not null;autoCreateTime" json:"created_at"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.UserID == uuid.Nil {
		u.UserID = uuid.New()
	}
	if u.Status == "" {
		u.Status = UserStatusActive
	}
	return nil
}
