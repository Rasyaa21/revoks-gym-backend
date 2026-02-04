package services

import (
	"fiber-gorm-app/internal/dto"
	"fiber-gorm-app/internal/repository"
	"time"
)

type NotificationService interface {
	MyList(userID uint, limit int) ([]dto.NotificationResponse, error)
	Detail(userID uint, id uint) (*dto.NotificationResponse, error)
	MarkRead(userID uint, id uint) error
}

type notificationService struct {
	repo repository.NotificationRepository
}

func NewNotificationService(repo repository.NotificationRepository) NotificationService {
	return &notificationService{repo: repo}
}

func (s *notificationService) MyList(userID uint, limit int) ([]dto.NotificationResponse, error) {
	items, err := s.repo.ListByUserID(userID, limit)
	if err != nil {
		return nil, err
	}
	resp := make([]dto.NotificationResponse, 0, len(items))
	for _, it := range items {
		resp = append(resp, dto.NotificationResponse{ID: it.ID, Title: it.Title, Body: it.Body, Type: it.Type, IsRead: it.IsRead, CreatedAt: it.CreatedAt.Format(time.RFC3339)})
	}
	return resp, nil
}

func (s *notificationService) Detail(userID uint, id uint) (*dto.NotificationResponse, error) {
	it, err := s.repo.FindByIDForUser(id, userID)
	if err != nil {
		return nil, err
	}
	return &dto.NotificationResponse{ID: it.ID, Title: it.Title, Body: it.Body, Type: it.Type, IsRead: it.IsRead, CreatedAt: it.CreatedAt.Format(time.RFC3339)}, nil
}

func (s *notificationService) MarkRead(userID uint, id uint) error {
	return s.repo.MarkRead(id, userID)
}
