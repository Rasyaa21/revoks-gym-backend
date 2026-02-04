package services

import (
	"fiber-gorm-app/internal/dto"
	"fiber-gorm-app/internal/repository"
	"time"
)

type TemplateService interface {
	MyFollowed(userID uint) ([]dto.TemplateResponse, error)
	Detail(id uint) (*dto.TemplateResponse, error)
	Follow(userID uint, templateID uint) error
}

type templateService struct {
	repo repository.TemplateRepository
}

func NewTemplateService(repo repository.TemplateRepository) TemplateService {
	return &templateService{repo: repo}
}

func (s *templateService) MyFollowed(userID uint) ([]dto.TemplateResponse, error) {
	items, err := s.repo.ListFollowedByUserID(userID)
	if err != nil {
		return nil, err
	}
	resp := make([]dto.TemplateResponse, 0, len(items))
	for _, it := range items {
		resp = append(resp, dto.TemplateResponse{ID: it.ID, Name: it.Name, Description: it.Description})
	}
	return resp, nil
}

func (s *templateService) Detail(id uint) (*dto.TemplateResponse, error) {
	t, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return &dto.TemplateResponse{ID: t.ID, Name: t.Name, Description: t.Description}, nil
}

func (s *templateService) Follow(userID uint, templateID uint) error {
	return s.repo.Follow(userID, templateID, time.Now())
}
