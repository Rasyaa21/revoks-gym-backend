package services

import (
	"errors"
	"fiber-gorm-app/internal/dto"
	"fiber-gorm-app/internal/models"
	"fiber-gorm-app/internal/repository"
	"time"
)

type TargetService interface {
	MyTargets(userID uint, period string) ([]dto.TargetResponse, error)
	CreateTarget(userID uint, req *dto.CreateTargetRequest) (*dto.TargetResponse, error)
	AddProgress(userID uint, targetID uint, req *dto.AddTargetProgressRequest) (*dto.TargetProgressResponse, error)
	ProgressHistory(userID uint, targetID uint, limit int) ([]dto.TargetProgressResponse, error)
}

type targetService struct {
	repo repository.TargetRepository
}

func NewTargetService(repo repository.TargetRepository) TargetService {
	return &targetService{repo: repo}
}

func (s *targetService) MyTargets(userID uint, period string) ([]dto.TargetResponse, error) {
	items, err := s.repo.ListByUserIDAndPeriod(userID, period)
	if err != nil {
		return nil, err
	}
	resp := make([]dto.TargetResponse, 0, len(items))
	for _, it := range items {
		resp = append(resp, dto.TargetResponse{
			ID:        it.ID,
			Period:    it.Period,
			Title:     it.Title,
			GoalValue: it.GoalValue,
			StartDate: it.StartDate.Format(time.RFC3339),
			EndDate:   it.EndDate.Format(time.RFC3339),
		})
	}
	return resp, nil
}

func (s *targetService) CreateTarget(userID uint, req *dto.CreateTargetRequest) (*dto.TargetResponse, error) {
	if req == nil {
		return nil, errors.New("invalid request")
	}
	if req.Period == "" {
		return nil, errors.New("period is required")
	}
	if req.Title == "" {
		return nil, errors.New("title is required")
	}
	if req.GoalValue <= 0 {
		return nil, errors.New("goal_value must be greater than 0")
	}
	start, err := time.Parse(time.RFC3339, req.StartDate)
	if err != nil {
		return nil, errors.New("start_date must be RFC3339")
	}
	end, err := time.Parse(time.RFC3339, req.EndDate)
	if err != nil {
		return nil, errors.New("end_date must be RFC3339")
	}
	if !end.After(start) {
		return nil, errors.New("end_date must be after start_date")
	}

	t := &models.Target{
		UserID:    userID,
		Period:    req.Period,
		Title:     req.Title,
		GoalValue: req.GoalValue,
		StartDate: start,
		EndDate:   end,
	}
	if err := s.repo.Create(t); err != nil {
		return nil, err
	}

	return &dto.TargetResponse{
		ID:        t.ID,
		Period:    t.Period,
		Title:     t.Title,
		GoalValue: t.GoalValue,
		StartDate: t.StartDate.Format(time.RFC3339),
		EndDate:   t.EndDate.Format(time.RFC3339),
	}, nil
}

func (s *targetService) AddProgress(userID uint, targetID uint, req *dto.AddTargetProgressRequest) (*dto.TargetProgressResponse, error) {
	target, err := s.repo.FindByID(targetID)
	if err != nil {
		return nil, err
	}
	if target.UserID != userID {
		return nil, errors.New("target not found")
	}

	when := time.Now()
	if req.RecordedAt != "" {
		if t, err := time.Parse(time.RFC3339, req.RecordedAt); err == nil {
			when = t
		}
	}
	p := &models.TargetProgress{TargetID: targetID, Value: req.Value, RecordedAt: when}
	if err := s.repo.CreateProgress(p); err != nil {
		return nil, err
	}
	return &dto.TargetProgressResponse{ID: p.ID, Value: p.Value, RecordedAt: p.RecordedAt.Format(time.RFC3339)}, nil
}

func (s *targetService) ProgressHistory(userID uint, targetID uint, limit int) ([]dto.TargetProgressResponse, error) {
	target, err := s.repo.FindByID(targetID)
	if err != nil {
		return nil, err
	}
	if target.UserID != userID {
		return nil, errors.New("target not found")
	}
	items, err := s.repo.ListProgress(targetID, limit)
	if err != nil {
		return nil, err
	}
	resp := make([]dto.TargetProgressResponse, 0, len(items))
	for _, it := range items {
		resp = append(resp, dto.TargetProgressResponse{ID: it.ID, Value: it.Value, RecordedAt: it.RecordedAt.Format(time.RFC3339)})
	}
	return resp, nil
}
