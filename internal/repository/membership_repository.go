package repository

import (
	"fiber-gorm-app/internal/models"
	"time"

	"gorm.io/gorm"
)

type MembershipRepository interface {
	Create(m *models.Membership) error
	ExpireActiveForUser(userID uint, expiredAt time.Time) error
	FindLatestByUserID(userID uint) (*models.Membership, error)
	ListByUserID(userID uint) ([]models.Membership, error)
}

type membershipRepository struct {
	db *gorm.DB
}

func NewMembershipRepository(db *gorm.DB) MembershipRepository {
	return &membershipRepository{db: db}
}

func (r *membershipRepository) Create(m *models.Membership) error {
	return r.db.Create(m).Error
}

func (r *membershipRepository) ExpireActiveForUser(userID uint, expiredAt time.Time) error {
	return r.db.Model(&models.Membership{}).
		Where("user_id = ? AND status = ?", userID, "active").
		Updates(map[string]interface{}{"status": "expired", "ends_at": expiredAt}).Error
}

func (r *membershipRepository) FindLatestByUserID(userID uint) (*models.Membership, error) {
	var m models.Membership
	if err := r.db.Where("user_id = ?", userID).Order("ends_at desc").First(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *membershipRepository) ListByUserID(userID uint) ([]models.Membership, error) {
	var ms []models.Membership
	if err := r.db.Where("user_id = ?", userID).Order("ends_at desc").Find(&ms).Error; err != nil {
		return nil, err
	}
	return ms, nil
}
