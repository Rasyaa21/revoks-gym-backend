package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Notification struct {
	NotificationID uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"notification_id"`
	UserID         uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	Type           string         `gorm:"type:varchar(50);not null" json:"type"`
	Title          string         `gorm:"type:varchar(255);not null" json:"title"`
	Body           *string        `gorm:"type:text" json:"body"`
	PayloadJSON    datatypes.JSON `gorm:"type:jsonb" json:"payload_json"`
	ReadAt         *time.Time     `gorm:"type:timestamp" json:"read_at"`
	CreatedAt      time.Time      `gorm:"type:timestamptz;not null;autoCreateTime" json:"created_at"`

	User *User `gorm:"foreignKey:UserID;references:UserID" json:"user,omitempty"`
}

func (Notification) TableName() string {
	return "notifications"
}

type AuditLog struct {
	AuditID     uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"audit_id"`
	ActorUserID *uuid.UUID     `gorm:"type:uuid" json:"actor_user_id"`
	Action      string         `gorm:"type:varchar(100);not null" json:"action"`
	EntityType  *string        `gorm:"type:varchar(100)" json:"entity_type"`
	EntityID    *uuid.UUID     `gorm:"type:uuid" json:"entity_id"`
	PayloadJSON datatypes.JSON `gorm:"type:jsonb" json:"payload_json"`
	CreatedAt   time.Time      `gorm:"type:timestamptz;not null;autoCreateTime" json:"created_at"`

	ActorUser *User `gorm:"foreignKey:ActorUserID;references:UserID" json:"actor_user,omitempty"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}
