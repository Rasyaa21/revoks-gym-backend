package repository

import (
	"fiber-gorm-app/internal/models"

	"gorm.io/gorm"
)

type TargetRepository interface {
	ListByUserIDAndPeriod(userID uint, period string) ([]models.Target, error)
	FindByID(id uint) (*models.Target, error)
	Create(t *models.Target) error
	CreateProgress(p *models.TargetProgress) error
	ListProgress(targetID uint, limit int) ([]models.TargetProgress, error)
}

type targetRepository struct {
	db *gorm.DB
}

func NewTargetRepository(db *gorm.DB) TargetRepository {
	return &targetRepository{db: db}
}

func (r *targetRepository) ListByUserIDAndPeriod(userID uint, period string) ([]models.Target, error) {
	var targets []models.Target
	q := r.db.Where("user_id = ?", userID).Order("start_date desc")
	if period != "" {
		q = q.Where("period = ?", period)
	}
	if err := q.Find(&targets).Error; err != nil {
		return nil, err
	}
	return targets, nil
}

func (r *targetRepository) FindByID(id uint) (*models.Target, error) {
	var t models.Target
	if err := r.db.First(&t, id).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *targetRepository) Create(t *models.Target) error {
	return r.db.Create(t).Error
}

func (r *targetRepository) CreateProgress(p *models.TargetProgress) error {
	return r.db.Create(p).Error
}

func (r *targetRepository) ListProgress(targetID uint, limit int) ([]models.TargetProgress, error) {
	var items []models.TargetProgress
	q := r.db.Where("target_id = ?", targetID).Order("recorded_at desc")
	if limit > 0 {
		q = q.Limit(limit)
	}
	if err := q.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
