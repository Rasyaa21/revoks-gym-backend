package repository

import (
	"fiber-gorm-app/internal/models"
	"time"

	"gorm.io/gorm"
)

type TemplateRepository interface {
	ListFollowedByUserID(userID uint) ([]models.Template, error)
	FindByID(id uint) (*models.Template, error)
	Follow(userID uint, templateID uint, followedAt time.Time) error
}

type templateRepository struct {
	db *gorm.DB
}

func NewTemplateRepository(db *gorm.DB) TemplateRepository {
	return &templateRepository{db: db}
}

func (r *templateRepository) ListFollowedByUserID(userID uint) ([]models.Template, error) {
	var templates []models.Template
	if err := r.db.Table("templates").
		Joins("join user_template_follows utf on utf.template_id = templates.id").
		Where("utf.user_id = ?", userID).
		Order("templates.name asc").
		Find(&templates).Error; err != nil {
		return nil, err
	}
	return templates, nil
}

func (r *templateRepository) FindByID(id uint) (*models.Template, error) {
	var t models.Template
	if err := r.db.First(&t, id).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *templateRepository) Follow(userID uint, templateID uint, followedAt time.Time) error {
	// idempotent-ish: try find first, otherwise create
	var existing models.UserTemplateFollow
	err := r.db.Where("user_id = ? AND template_id = ?", userID, templateID).First(&existing).Error
	if err == nil {
		return nil
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	return r.db.Create(&models.UserTemplateFollow{UserID: userID, TemplateID: templateID, FollowedAt: followedAt}).Error
}
