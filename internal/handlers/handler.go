package handlers

import "fiber-gorm-app/internal/services"

// Handler holds all handler instances
type Handler struct {
	User         *UserHandler
	Auth         *AuthHandler
	Me           *MeHandler
	Membership   *MembershipHandler
	QR           *QRHandler
	Attendance   *AttendanceHandler
	Workout      *WorkoutHandler
	Template     *TemplateHandler
	Target       *TargetHandler
	Trainer      *TrainerHandler
	Notification *NotificationHandler
	Setting      *SettingHandler
}

// NewHandler creates a new handler with all dependencies
func NewHandler(svc *services.Service) *Handler {
	return &Handler{
		User:         NewUserHandler(svc.User),
		Auth:         NewAuthHandler(svc.Auth),
		Me:           NewMeHandler(svc.User),
		Membership:   NewMembershipHandler(svc.Membership),
		QR:           NewQRHandler(svc.QR),
		Attendance:   NewAttendanceHandler(svc.Attendance),
		Workout:      NewWorkoutHandler(svc.Workout),
		Template:     NewTemplateHandler(svc.Template),
		Target:       NewTargetHandler(svc.Target),
		Trainer:      NewTrainerHandler(svc.Trainer),
		Notification: NewNotificationHandler(svc.Notification),
		Setting:      NewSettingHandler(svc.Setting),
	}
}
