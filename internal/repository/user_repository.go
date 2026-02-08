package repository

import (
	"fiber-gorm-app/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error
	FindAll(page, perPage int) ([]models.User, int64, error)
	FindByID(id uuid.UUID) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	Update(user *models.User) error
	Delete(id uuid.UUID) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindAll(page, perPage int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	offset := (page - 1) * perPage

	if err := r.db.Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Offset(offset).Limit(perPage).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *userRepository) FindByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := r.db.Where("user_id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id uuid.UUID) error {
	return r.db.Where("user_id = ?", id).Delete(&models.User{}).Error
}
