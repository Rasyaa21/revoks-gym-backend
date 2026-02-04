package handlers

import "fiber-gorm-app/internal/services"

// Handler holds all handler instances
type Handler struct {
	User *UserHandler
}

// NewHandler creates a new handler with all dependencies
func NewHandler(svc *services.Service) *Handler {
	return &Handler{
		User: NewUserHandler(svc.User),
	}
}
