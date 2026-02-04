package repository

import (
	"fiber-gorm-app/internal/models"

	"gorm.io/gorm"
)

type AttendanceRepository interface {
	Create(log *models.AttendanceLog) error
	FindLastByUserID(userID uint) (*models.AttendanceLog, error)
	ListByUserID(userID uint, limit int) ([]models.AttendanceLog, error)
}

type attendanceRepository struct {
	db *gorm.DB
}

func NewAttendanceRepository(db *gorm.DB) AttendanceRepository {
	return &attendanceRepository{db: db}
}

func (r *attendanceRepository) Create(log *models.AttendanceLog) error {
	return r.db.Create(log).Error
}

func (r *attendanceRepository) FindLastByUserID(userID uint) (*models.AttendanceLog, error) {
	var log models.AttendanceLog
	if err := r.db.Where("user_id = ?", userID).Order("occurred_at desc").First(&log).Error; err != nil {
		return nil, err
	}
	return &log, nil
}

func (r *attendanceRepository) ListByUserID(userID uint, limit int) ([]models.AttendanceLog, error) {
	var logs []models.AttendanceLog
	q := r.db.Where("user_id = ?", userID).Order("occurred_at desc")
	if limit > 0 {
		q = q.Limit(limit)
	}
	if err := q.Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}
