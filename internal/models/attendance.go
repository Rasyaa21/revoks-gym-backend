package models

import (
	"time"

	"github.com/google/uuid"
)

type ScannerDevice struct {
	ScannerID  uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"scanner_id"`
	Name       string    `gorm:"type:varchar(100);not null" json:"name"`
	Location   *string   `gorm:"type:varchar(200)" json:"location"`
	APIKeyHash string    `gorm:"type:varchar(255);not null" json:"api_key_hash"`
	IsActive   bool      `gorm:"not null;default:true" json:"is_active"`
	CreatedAt  time.Time `gorm:"type:timestamptz;not null;autoCreateTime" json:"created_at"`
}

func (ScannerDevice) TableName() string {
	return "scanner_devices"
}

type UserAccessState struct {
	UserID        uuid.UUID  `gorm:"type:uuid;primaryKey" json:"user_id"`
	GateState     string     `gorm:"type:varchar(20);not null;default:'OUTSIDE'" json:"gate_state"`
	LastScannedAt *time.Time `gorm:"type:timestamp" json:"last_scanned_at"`
	LastScannerID *uuid.UUID `gorm:"type:uuid" json:"last_scanner_id"`

	// Relationships
	User    *User          `gorm:"foreignKey:UserID;references:UserID" json:"user,omitempty"`
	Scanner *ScannerDevice `gorm:"foreignKey:LastScannerID;references:ScannerID" json:"scanner,omitempty"`
}

func (UserAccessState) TableName() string {
	return "user_access_state"
}

type AttendanceLog struct {
	AttendanceLogID  uuid.UUID           `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"attendance_log_id"`
	ScannedAt        time.Time           `gorm:"type:timestamptz;not null;autoCreateTime" json:"scanned_at"`
	ScannerID        uuid.UUID           `gorm:"type:uuid;not null" json:"scanner_id"`
	UserID           *uuid.UUID          `gorm:"type:uuid" json:"user_id"`
	OneTimePassID    *uuid.UUID          `gorm:"type:uuid" json:"one_time_pass_id"`
	ScannedTokenHash *string             `gorm:"type:varchar(255)" json:"scanned_token_hash"`
	Direction        *AttendanceDirection `gorm:"type:varchar(10)" json:"direction"`
	Result           AttendanceResult    `gorm:"type:varchar(20);not null" json:"result"`
	RejectReason     *string             `gorm:"type:varchar(255)" json:"reject_reason"`

	// Relationships
	Scanner     *ScannerDevice `gorm:"foreignKey:ScannerID;references:ScannerID" json:"scanner,omitempty"`
	User        *User          `gorm:"foreignKey:UserID;references:UserID" json:"user,omitempty"`
	OneTimePass *OneTimePass   `gorm:"foreignKey:OneTimePassID;references:OneTimePassID" json:"one_time_pass,omitempty"`
}

func (AttendanceLog) TableName() string {
	return "attendance_logs"
}
