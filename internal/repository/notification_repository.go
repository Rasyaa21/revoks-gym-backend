package repository

import (
	"fiber-gorm-app/internal/models"

	"gorm.io/gorm"
)

type NotificationRepository interface {
	ListByUserID(userID uint, limit int) ([]models.Notification, error)
	FindByIDForUser(id uint, userID uint) (*models.Notification, error)
	MarkRead(id uint, userID uint) error
}

type notificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) ListByUserID(userID uint, limit int) ([]models.Notification, error) {
	var items []models.Notification
	q := r.db.Where("user_id = ?", userID).Order("created_at desc")
	if limit > 0 {
		q = q.Limit(limit)
	}
	if err := q.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *notificationRepository) FindByIDForUser(id uint, userID uint) (*models.Notification, error) {
	var n models.Notification
	if err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&n).Error; err != nil {
		return nil, err
	}
	return &n, nil
}

func (r *notificationRepository) MarkRead(id uint, userID uint) error {
	return r.db.Model(&models.Notification{}).Where("id = ? AND user_id = ?", id, userID).Update("is_read", true).Error
}
