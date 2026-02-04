package routes

import (
	"fiber-gorm-app/internal/config"
	"fiber-gorm-app/internal/handlers"
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

	// User routes
	users := v1.Group("/users")
	users.Post("/", handler.User.CreateUser)
	users.Get("/", handler.User.GetAllUsers)
	users.Get("/:id", handler.User.GetUserByID)
	users.Put("/:id", handler.User.UpdateUser)
	users.Delete("/:id", handler.User.DeleteUser)
	users.Put("/:id/password", handler.User.ChangePassword)
}
