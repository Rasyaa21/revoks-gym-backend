package services

import (
	"fiber-gorm-app/internal/dto"
	"fiber-gorm-app/internal/repository"
)

type TrainerService interface {
	List() ([]dto.TrainerResponse, error)
	Detail(id uint) (*dto.TrainerResponse, error)
	Schedule(id uint) ([]dto.TrainerScheduleResponse, error)
}

type trainerService struct {
	repo repository.TrainerRepository
}

func NewTrainerService(repo repository.TrainerRepository) TrainerService {
	return &trainerService{repo: repo}
}

func (s *trainerService) List() ([]dto.TrainerResponse, error) {
	items, err := s.repo.List()
	if err != nil {
		return nil, err
	}
	resp := make([]dto.TrainerResponse, 0, len(items))
	for _, it := range items {
		resp = append(resp, dto.TrainerResponse{ID: it.ID, Name: it.Name, Bio: it.Bio, Specialty: it.Specialty, PhotoURL: it.PhotoURL})
	}
	return resp, nil
}

func (s *trainerService) Detail(id uint) (*dto.TrainerResponse, error) {
	it, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return &dto.TrainerResponse{ID: it.ID, Name: it.Name, Bio: it.Bio, Specialty: it.Specialty, PhotoURL: it.PhotoURL}, nil
}

func (s *trainerService) Schedule(id uint) ([]dto.TrainerScheduleResponse, error) {
	items, err := s.repo.ListSchedule(id)
	if err != nil {
		return nil, err
	}
	resp := make([]dto.TrainerScheduleResponse, 0, len(items))
	for _, it := range items {
		resp = append(resp, dto.TrainerScheduleResponse{ID: it.ID, DayOfWeek: it.DayOfWeek, StartTime: it.StartTime, EndTime: it.EndTime, Location: it.Location})
	}
	return resp, nil
}
