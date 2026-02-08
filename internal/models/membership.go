package models

import (
	"time"

	"github.com/google/uuid"
)

type MembershipPlan struct {
	PlanID       uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"plan_id"`
	Name         string    `gorm:"type:varchar(100);not null" json:"name"`
	PriceAmount  float64   `gorm:"type:decimal;not null" json:"price_amount"`
	DurationDays *int      `gorm:"type:int" json:"duration_days"`
	SessionQuota *int      `gorm:"type:int" json:"session_quota"`
	IsActive     bool      `gorm:"not null;default:true" json:"is_active"`
	CreatedAt    time.Time `gorm:"type:timestamptz;not null;autoCreateTime" json:"created_at"`
}

func (MembershipPlan) TableName() string {
	return "membership_plans"
}

type MemberMembership struct {
	MembershipID       uuid.UUID        `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"membership_id"`
	MemberID           uuid.UUID        `gorm:"type:uuid;not null" json:"member_id"`
	PlanID             uuid.UUID        `gorm:"type:uuid;not null" json:"plan_id"`
	StartAt            *time.Time       `gorm:"type:timestamp" json:"start_at"`
	ExpireAt           *time.Time       `gorm:"type:timestamp" json:"expire_at"`
	Status             MembershipStatus `gorm:"type:varchar(20);not null" json:"status"`
	TotalSessions      *int             `gorm:"type:int" json:"total_sessions"`
	UsedSessions       *int             `gorm:"type:int" json:"used_sessions"`
	RemainingSessions  *int             `gorm:"type:int" json:"remaining_sessions"`
	ActivatedByOrderID *uuid.UUID       `gorm:"type:uuid" json:"activated_by_order_id"`
	CreatedAt          time.Time        `gorm:"type:timestamptz;not null;autoCreateTime" json:"created_at"`

	Member           *Member         `gorm:"foreignKey:MemberID;references:MemberID" json:"member,omitempty"`
	Plan             *MembershipPlan `gorm:"foreignKey:PlanID;references:PlanID" json:"plan,omitempty"`
	ActivatedByOrder *Order          `gorm:"foreignKey:ActivatedByOrderID;references:OrderID" json:"activated_by_order,omitempty"`
}

func (MemberMembership) TableName() string {
	return "member_memberships"
}

type SessionUsageLog struct {
	UsageID      uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"usage_id"`
	MembershipID uuid.UUID `gorm:"type:uuid;not null" json:"membership_id"`
	UsedAt       time.Time `gorm:"type:timestamptz;not null;autoCreateTime" json:"used_at"`
	UsedBy       *string   `gorm:"type:varchar(100)" json:"used_by"`
	Notes        *string   `gorm:"type:varchar(500)" json:"notes"`

	Membership *MemberMembership `gorm:"foreignKey:MembershipID;references:MembershipID" json:"membership,omitempty"`
}

func (SessionUsageLog) TableName() string {
	return "session_usage_logs"
}
