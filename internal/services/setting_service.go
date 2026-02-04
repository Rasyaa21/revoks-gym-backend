package services

import (
	"fiber-gorm-app/internal/dto"
	"fiber-gorm-app/internal/repository"
)

type SettingService interface {
	Get(userID uint) (*dto.SettingsResponse, error)
	Update(userID uint, req *dto.UpdateSettingsRequest) (*dto.SettingsResponse, error)
}

type settingService struct {
	repo repository.SettingRepository
}

func NewSettingService(repo repository.SettingRepository) SettingService {
	return &settingService{repo: repo}
}

func (s *settingService) Get(userID uint) (*dto.SettingsResponse, error) {
	setting, err := s.repo.GetOrCreateByUserID(userID)
	if err != nil {
		return nil, err
	}
	return &dto.SettingsResponse{PushEnabled: setting.PushEnabled, EmailEnabled: setting.EmailEnabled, Language: setting.Language}, nil
}

func (s *settingService) Update(userID uint, req *dto.UpdateSettingsRequest) (*dto.SettingsResponse, error) {
	setting, err := s.repo.GetOrCreateByUserID(userID)
	if err != nil {
		return nil, err
	}

	if req.PushEnabled != nil {
		setting.PushEnabled = *req.PushEnabled
	}
	if req.EmailEnabled != nil {
		setting.EmailEnabled = *req.EmailEnabled
	}
	if req.Language != "" {
		setting.Language = req.Language
	}

	if err := s.repo.Update(setting); err != nil {
		return nil, err
	}
	return &dto.SettingsResponse{PushEnabled: setting.PushEnabled, EmailEnabled: setting.EmailEnabled, Language: setting.Language}, nil
}
