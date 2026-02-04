package repository

import "gorm.io/gorm"

// Repository holds all repository instances
type Repository struct {
	User UserRepository
}

// NewRepository creates a new repository with all dependencies
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		User: NewUserRepository(db),
	}
}
