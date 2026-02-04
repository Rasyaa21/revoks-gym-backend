package services

import (
	"fiber-gorm-app/internal/dto"
	"fiber-gorm-app/internal/repository"
	"time"
)

type AttendanceService interface {
	MyHistory(userID uint, limit int) ([]dto.AttendanceLogResponse, error)
}

type attendanceService struct {
	repo repository.AttendanceRepository
}

func NewAttendanceService(repo repository.AttendanceRepository) AttendanceService {
	return &attendanceService{repo: repo}
}

func (s *attendanceService) MyHistory(userID uint, limit int) ([]dto.AttendanceLogResponse, error) {
	logs, err := s.repo.ListByUserID(userID, limit)
	if err != nil {
		return nil, err
	}

	resp := make([]dto.AttendanceLogResponse, 0, len(logs))
	for _, l := range logs {
		resp = append(resp, dto.AttendanceLogResponse{
			ID:         l.ID,
			Direction:  l.Direction,
			Source:     l.Source,
			OccurredAt: l.OccurredAt.Format(time.RFC3339),
		})
	}
	return resp, nil
}
