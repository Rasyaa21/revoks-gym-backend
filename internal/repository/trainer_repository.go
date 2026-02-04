package repository

import (
	"fiber-gorm-app/internal/models"

	"gorm.io/gorm"
)

type TrainerRepository interface {
	List() ([]models.Trainer, error)
	FindByID(id uint) (*models.Trainer, error)
	ListSchedule(trainerID uint) ([]models.TrainerSchedule, error)
}

type trainerRepository struct {
	db *gorm.DB
}

func NewTrainerRepository(db *gorm.DB) TrainerRepository {
	return &trainerRepository{db: db}
}

func (r *trainerRepository) List() ([]models.Trainer, error) {
	var trainers []models.Trainer
	if err := r.db.Order("name asc").Find(&trainers).Error; err != nil {
		return nil, err
	}
	return trainers, nil
}

func (r *trainerRepository) FindByID(id uint) (*models.Trainer, error) {
	var t models.Trainer
	if err := r.db.First(&t, id).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *trainerRepository) ListSchedule(trainerID uint) ([]models.TrainerSchedule, error) {
	var items []models.TrainerSchedule
	if err := r.db.Where("trainer_id = ?", trainerID).Order("day_of_week asc, start_time asc").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
