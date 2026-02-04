package repository

import (
	"fiber-gorm-app/internal/models"

	"gorm.io/gorm"
)

type WorkoutRepository interface {
	Create(progress *models.WorkoutProgress) error
	ListByUserID(userID uint, limit int) ([]models.WorkoutProgress, error)
}

type workoutRepository struct {
	db *gorm.DB
}

func NewWorkoutRepository(db *gorm.DB) WorkoutRepository {
	return &workoutRepository{db: db}
}

func (r *workoutRepository) Create(progress *models.WorkoutProgress) error {
	return r.db.Create(progress).Error
}

func (r *workoutRepository) ListByUserID(userID uint, limit int) ([]models.WorkoutProgress, error) {
	var items []models.WorkoutProgress
	q := r.db.Where("user_id = ?", userID).Order("recorded_at desc")
	if limit > 0 {
		q = q.Limit(limit)
	}
	if err := q.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
