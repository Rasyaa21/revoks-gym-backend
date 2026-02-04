package repository

import "gorm.io/gorm"

// Repository holds all repository instances
type Repository struct {
	User         UserRepository
	Membership   MembershipRepository
	Attendance   AttendanceRepository
	Workout      WorkoutRepository
	Template     TemplateRepository
	Target       TargetRepository
	Trainer      TrainerRepository
	Notification NotificationRepository
	Setting      SettingRepository
}

// NewRepository creates a new repository with all dependencies
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		User:         NewUserRepository(db),
		Membership:   NewMembershipRepository(db),
		Attendance:   NewAttendanceRepository(db),
		Workout:      NewWorkoutRepository(db),
		Template:     NewTemplateRepository(db),
		Target:       NewTargetRepository(db),
		Trainer:      NewTrainerRepository(db),
		Notification: NewNotificationRepository(db),
		Setting:      NewSettingRepository(db),
	}
}
