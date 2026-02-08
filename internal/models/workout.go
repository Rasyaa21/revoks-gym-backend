package models

import (
	"time"

	"github.com/google/uuid"
)

type ExerciseCatalog struct {
	ExerciseID    uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"exercise_id"`
	Name          string    `gorm:"type:varchar(200);not null" json:"name"`
	PrimaryMuscle *string   `gorm:"type:varchar(100)" json:"primary_muscle"`
	Equipment     *string   `gorm:"type:varchar(100)" json:"equipment"`
	VideoURL      *string   `gorm:"type:varchar(500)" json:"video_url"`
	ImageURL      *string   `gorm:"type:varchar(500)" json:"image_url"`
	IsActive      bool      `gorm:"not null;default:true" json:"is_active"`
	CreatedAt     time.Time `gorm:"type:timestamptz;not null;autoCreateTime" json:"created_at"`
}

func (ExerciseCatalog) TableName() string {
	return "exercise_catalog"
}

type WorkoutTemplate struct {
	TemplateID      uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"template_id"`
	Title           string     `gorm:"type:varchar(200);not null" json:"title"`
	Description     *string    `gorm:"type:text" json:"description"`
	FocusMuscle     *string    `gorm:"type:varchar(100)" json:"focus_muscle"`
	Level           *string    `gorm:"type:varchar(50)" json:"level"`
	CreatedByUserID *uuid.UUID `gorm:"type:uuid" json:"created_by_user_id"`
	CoachPTID       *uuid.UUID `gorm:"type:uuid" json:"coach_pt_id"`
	IsActive        bool       `gorm:"not null;default:true" json:"is_active"`
	CreatedAt       time.Time  `gorm:"type:timestamptz;not null;autoCreateTime" json:"created_at"`

	CreatedByUser *User            `gorm:"foreignKey:CreatedByUserID;references:UserID" json:"created_by_user,omitempty"`
	CoachPT       *PersonalTrainer `gorm:"foreignKey:CoachPTID;references:PTID" json:"coach_pt,omitempty"`
}

func (WorkoutTemplate) TableName() string {
	return "workout_templates"
}

type WorkoutTemplateItem struct {
	ItemID      uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"item_id"`
	TemplateID  uuid.UUID `gorm:"type:uuid;not null" json:"template_id"`
	ExerciseID  uuid.UUID `gorm:"type:uuid;not null" json:"exercise_id"`
	OrderNo     *int      `gorm:"type:int" json:"order_no"`
	Sets        *int      `gorm:"type:int" json:"sets"`
	Reps        *string   `gorm:"type:varchar(50)" json:"reps"`
	RestSeconds *int      `gorm:"type:int" json:"rest_seconds"`
	Tempo       *string   `gorm:"type:varchar(20)" json:"tempo"`
	RPE         *int      `gorm:"type:int" json:"rpe"`
	Notes       *string   `gorm:"type:text" json:"notes"`

	Template *WorkoutTemplate `gorm:"foreignKey:TemplateID;references:TemplateID" json:"template,omitempty"`
	Exercise *ExerciseCatalog `gorm:"foreignKey:ExerciseID;references:ExerciseID" json:"exercise,omitempty"`
}

func (WorkoutTemplateItem) TableName() string {
	return "workout_template_items"
}

type MemberWorkoutFollow struct {
	FollowID   uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"follow_id"`
	MemberID   uuid.UUID `gorm:"type:uuid;not null" json:"member_id"`
	TemplateID uuid.UUID `gorm:"type:uuid;not null" json:"template_id"`
	FollowedAt time.Time `gorm:"type:timestamptz;not null;autoCreateTime" json:"followed_at"`
	Status     string    `gorm:"type:varchar(20);not null;default:'ACTIVE'" json:"status"`

	Member   *Member          `gorm:"foreignKey:MemberID;references:MemberID" json:"member,omitempty"`
	Template *WorkoutTemplate `gorm:"foreignKey:TemplateID;references:TemplateID" json:"template,omitempty"`
}

func (MemberWorkoutFollow) TableName() string {
	return "member_workout_follows"
}

type MemberWorkoutProgress struct {
	ProgressID  uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"progress_id"`
	FollowID    uuid.UUID `gorm:"type:uuid;not null" json:"follow_id"`
	CompletedAt time.Time `gorm:"type:timestamptz;not null;autoCreateTime" json:"completed_at"`
	Notes       *string   `gorm:"type:text" json:"notes"`

	Follow *MemberWorkoutFollow `gorm:"foreignKey:FollowID;references:FollowID" json:"follow,omitempty"`
}

func (MemberWorkoutProgress) TableName() string {
	return "member_workout_progress"
}
