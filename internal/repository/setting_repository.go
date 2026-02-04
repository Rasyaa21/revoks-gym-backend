package repository

import (
	"fiber-gorm-app/internal/models"

	"gorm.io/gorm"
)

type SettingRepository interface {
	GetOrCreateByUserID(userID uint) (*models.UserSetting, error)
	Update(setting *models.UserSetting) error
}

type settingRepository struct {
	db *gorm.DB
}

func NewSettingRepository(db *gorm.DB) SettingRepository {
	return &settingRepository{db: db}
}

func (r *settingRepository) GetOrCreateByUserID(userID uint) (*models.UserSetting, error) {
	var s models.UserSetting
	err := r.db.Where("user_id = ?", userID).First(&s).Error
	if err == nil {
		return &s, nil
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	// create default
	s = models.UserSetting{UserID: userID, PushEnabled: true, EmailEnabled: true, Language: "id"}
	if err := r.db.Create(&s).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *settingRepository) Update(setting *models.UserSetting) error {
	return r.db.Save(setting).Error
}
