package services

import "fiber-gorm-app/internal/repository"

// Service holds all service instances
type Service struct {
	User UserService
}

// NewService creates a new service with all dependencies
func NewService(repo *repository.Repository) *Service {
	return &Service{
		User: NewUserService(repo.User),
	}
}
