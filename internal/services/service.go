package services

import "fiber-gorm-app/internal/repository"

// Service holds all service instances
type Service struct {
	User         UserService
	Auth         AuthService
	Membership   MembershipService
	QR           QRService
	Attendance   AttendanceService
	Workout      WorkoutService
	Template     TemplateService
	Target       TargetService
	Trainer      TrainerService
	Notification NotificationService
	Setting      SettingService
}

// NewService creates a new service with all dependencies
func NewService(repo *repository.Repository) *Service {
	return &Service{
		User:         NewUserService(repo.User),
		Auth:         NewAuthService(repo.User),
		Membership:   NewMembershipService(repo.Membership),
		QR:           NewQRService(repo.Membership, repo.Attendance),
		Attendance:   NewAttendanceService(repo.Attendance),
		Workout:      NewWorkoutService(repo.Workout),
		Template:     NewTemplateService(repo.Template),
		Target:       NewTargetService(repo.Target),
		Trainer:      NewTrainerService(repo.Trainer),
		Notification: NewNotificationService(repo.Notification),
		Setting:      NewSettingService(repo.Setting),
	}
}
