package models

import (
	"time"

	"github.com/google/uuid"
)

type Member struct {
	MemberID        uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"member_id"`
	UserID          uuid.UUID  `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`
	LegacyMemberNo  *int64     `gorm:"type:bigint" json:"legacy_member_no"`
	LegacyRecordID  *int64     `gorm:"type:bigint" json:"legacy_record_id"`
	Gender          *string    `gorm:"type:varchar(10)" json:"gender"`
	BirthDate       *time.Time `gorm:"type:date" json:"birth_date"`
	Address         *string    `gorm:"type:text" json:"address"`
	PhoneOverride   *string    `gorm:"type:varchar(50)" json:"phone_override"`
	QRStaticHash    string     `gorm:"type:varchar(255);not null;uniqueIndex" json:"qr_static_hash"`
	CreatedAt       time.Time  `gorm:"type:timestamptz;not null;autoCreateTime" json:"created_at"`

	// Relationships (single references only, no slices)
	User *User `gorm:"foreignKey:UserID;references:UserID" json:"user,omitempty"`
}

func (Member) TableName() string {
	return "members"
}

type PersonalTrainer struct {
	PTID            uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"pt_id"`
	UserID          uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`
	Bio             *string   `gorm:"type:text" json:"bio"`
	ExperienceYears *int      `gorm:"type:int" json:"experience_years"`
	Specialties     *string   `gorm:"type:text" json:"specialties"`
	BaseRateAmount  *float64  `gorm:"type:decimal" json:"base_rate_amount"`
	RatingAvg       *float64  `gorm:"type:decimal" json:"rating_avg"`
	RatingCount     *int      `gorm:"type:int" json:"rating_count"`
	IsActive        bool      `gorm:"not null;default:true" json:"is_active"`
	CreatedAt       time.Time `gorm:"type:timestamptz;not null;autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time `gorm:"type:timestamptz;not null;autoUpdateTime" json:"updated_at"`

	// Relationships (single references only, no slices)
	User *User `gorm:"foreignKey:UserID;references:UserID" json:"user,omitempty"`
}

func (PersonalTrainer) TableName() string {
	return "personal_trainers"
}

type PTReview struct {
	ReviewID  uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"review_id"`
	PTID      uuid.UUID `gorm:"type:uuid;not null" json:"pt_id"`
	MemberID  uuid.UUID `gorm:"type:uuid;not null" json:"member_id"`
	Rating    *int      `gorm:"type:int" json:"rating"`
	Comment   *string   `gorm:"type:text" json:"comment"`
	CreatedAt time.Time `gorm:"type:timestamptz;not null;autoCreateTime" json:"created_at"`

	// Relationships
	PersonalTrainer *PersonalTrainer `gorm:"foreignKey:PTID;references:PTID" json:"personal_trainer,omitempty"`
	Member          *Member          `gorm:"foreignKey:MemberID;references:MemberID" json:"member,omitempty"`
}

func (PTReview) TableName() string {
	return "pt_reviews"
}

type PTAvailabilityRule struct {
	RuleID    uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"rule_id"`
	PTID      uuid.UUID `gorm:"type:uuid;not null" json:"pt_id"`
	Weekday   string    `gorm:"type:varchar(10);not null" json:"weekday"`
	StartTime string    `gorm:"type:time;not null" json:"start_time"`
	EndTime   string    `gorm:"type:time;not null" json:"end_time"`
	IsActive  bool      `gorm:"not null;default:true" json:"is_active"`
	CreatedAt time.Time `gorm:"type:timestamptz;not null;autoCreateTime" json:"created_at"`

	// Relationships
	PersonalTrainer *PersonalTrainer `gorm:"foreignKey:PTID;references:PTID" json:"personal_trainer,omitempty"`
}

func (PTAvailabilityRule) TableName() string {
	return "pt_availability_rules"
}
