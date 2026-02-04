package services

import (
	"fiber-gorm-app/internal/dto"
	"fiber-gorm-app/internal/models"
	"fiber-gorm-app/internal/repository"
	"time"
)

type WorkoutService interface {
	Create(userID uint, req *dto.CreateWorkoutProgressRequest) (*dto.WorkoutProgressResponse, error)
	MyHistory(userID uint, limit int) ([]dto.WorkoutProgressResponse, error)
}

type workoutService struct {
	repo repository.WorkoutRepository
}

func NewWorkoutService(repo repository.WorkoutRepository) WorkoutService {
	return &workoutService{repo: repo}
}

func (s *workoutService) Create(userID uint, req *dto.CreateWorkoutProgressRequest) (*dto.WorkoutProgressResponse, error) {
	recordedAt := time.Now()
	if req.RecordedAt != "" {
		if t, err := time.Parse(time.RFC3339, req.RecordedAt); err == nil {
			recordedAt = t
		}
	}
	item := &models.WorkoutProgress{
		UserID:     userID,
		Title:      req.Title,
		Notes:      req.Notes,
		RecordedAt: recordedAt,
	}
	if err := s.repo.Create(item); err != nil {
		return nil, err
	}
	return &dto.WorkoutProgressResponse{
		ID:         item.ID,
		Title:      item.Title,
		Notes:      item.Notes,
		RecordedAt: item.RecordedAt.Format(time.RFC3339),
	}, nil
}

func (s *workoutService) MyHistory(userID uint, limit int) ([]dto.WorkoutProgressResponse, error) {
	items, err := s.repo.ListByUserID(userID, limit)
	if err != nil {
		return nil, err
	}
	resp := make([]dto.WorkoutProgressResponse, 0, len(items))
	for _, it := range items {
		resp = append(resp, dto.WorkoutProgressResponse{
			ID:         it.ID,
			Title:      it.Title,
			Notes:      it.Notes,
			RecordedAt: it.RecordedAt.Format(time.RFC3339),
		})
	}
	return resp, nil
}
