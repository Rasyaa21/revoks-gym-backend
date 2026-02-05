package routes

import (
	"fiber-gorm-app/internal/config"
	"fiber-gorm-app/internal/handlers"
	"fiber-gorm-app/internal/middleware"
	"fiber-gorm-app/internal/repository"
	"fiber-gorm-app/internal/services"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	// Initialize dependencies
	db := config.GetDB()
	repo := repository.NewRepository(db)
	svc := services.NewService(repo)
	handler := handlers.NewHandler(svc)

	// API routes
	api := app.Group("/api")
	v1 := api.Group("/v1")

	// Auth routes
	auth := v1.Group("/auth")
	auth.Post("/login", handler.Auth.Login)

	// Public QR scan endpoint (for gate device)
	qrPublic := v1.Group("/qr")
	qrPublic.Post("/scan", handler.QR.Scan)

	// User routes
	users := v1.Group("/users")
	users.Post("/", handler.User.CreateUser)
	users.Get("/", handler.User.GetAllUsers)
	users.Get("/:id", handler.User.GetUserByID)
	users.Put("/:id", handler.User.UpdateUser)
	users.Delete("/:id", handler.User.DeleteUser)
	users.Put("/:id/password", handler.User.ChangePassword)

	// Authenticated routes
	authed := v1.Group("", middleware.RequireAuth())
	authed.Get("/me", handler.Me.GetMe)

	// Membership
	authed.Get("/membership", handler.Membership.GetMyMembership)
	authed.Post("/membership/renew", handler.Membership.Renew)

	// QR access (generate)
	authed.Get("/qr/code", handler.QR.GetMyQRCode)

	// Notifications (global)
	authed.Get("/notifications", handler.Notification.MyList)
	authed.Get("/notifications/:id", handler.Notification.Detail)
	authed.Put("/notifications/:id/read", handler.Notification.MarkRead)

	// Activity
	authed.Get("/attendance/history", handler.Attendance.MyHistory)
	authed.Get("/workouts/progress", handler.Workout.MyHistory)
	authed.Post("/workouts/progress", handler.Workout.Create)

	// Program
	authed.Get("/templates/followed", handler.Template.MyFollowed)
	authed.Get("/templates/:id", handler.Template.Detail)
	authed.Post("/templates/follow", handler.Template.Follow)
	authed.Get("/targets", handler.Target.MyTargets)
	authed.Post("/targets", handler.Target.Create)
	authed.Get("/targets/:id/progress", handler.Target.ProgressHistory)
	authed.Post("/targets/:id/progress", handler.Target.AddProgress)

	// PT
	authed.Get("/pt", handler.Trainer.List)
	authed.Get("/pt/:id", handler.Trainer.Detail)
	authed.Get("/pt/:id/schedule", handler.Trainer.Schedule)

	// Settings
	authed.Get("/settings", handler.Setting.Get)
	authed.Put("/settings", handler.Setting.Update)
}
