package services

import (
	"fiber-gorm-app/internal/dto"
	"fiber-gorm-app/internal/models"
	"fiber-gorm-app/internal/repository"
	"fiber-gorm-app/internal/utils"
	"time"

	"gorm.io/gorm"
)

type QRService interface {
	GetMyQRCode(userID uint) (*dto.MyQRCodeResponse, error)
	Scan(token string) (*dto.ScanQRResponse, error)
}

type qrService struct {
	membershipRepo repository.MembershipRepository
	attendanceRepo repository.AttendanceRepository
}

func NewQRService(membershipRepo repository.MembershipRepository, attendanceRepo repository.AttendanceRepository) QRService {
	return &qrService{membershipRepo: membershipRepo, attendanceRepo: attendanceRepo}
}

func (s *qrService) GetMyQRCode(userID uint) (*dto.MyQRCodeResponse, error) {
	status := "expired"
	if current, err := s.membershipRepo.FindLatestByUserID(userID); err == nil {
		if current.Status == "active" && time.Now().Before(current.EndsAt) {
			status = "active"
		}
	}

	ttl := 60 * time.Second
	token, err := utils.NewQRToken(userID, ttl)
	if err != nil {
		return nil, err
	}

	return &dto.MyQRCodeResponse{
		Token:            token,
		ExpiresAt:        time.Now().Add(ttl).Format(time.RFC3339),
		MembershipStatus: status,
	}, nil
}

func (s *qrService) Scan(token string) (*dto.ScanQRResponse, error) {
	claims, err := utils.ParseToken(token, "qr")
	if err != nil {
		return &dto.ScanQRResponse{Accepted: false, Reason: "invalid_qr"}, nil
	}

	// membership must be active
	current, err := s.membershipRepo.FindLatestByUserID(claims.UserID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &dto.ScanQRResponse{Accepted: false, Reason: "membership_required"}, nil
		}
		return nil, err
	}
	if current.Status != "active" || time.Now().After(current.EndsAt) {
		return &dto.ScanQRResponse{Accepted: false, Reason: "membership_expired"}, nil
	}

	// toggle IN/OUT based on last
	nextDirection := "in"
	if last, err := s.attendanceRepo.FindLastByUserID(claims.UserID); err == nil {
		if last.Direction == "in" {
			nextDirection = "out"
		}
	}

	now := time.Now()
	log := &models.AttendanceLog{
		UserID:     claims.UserID,
		Direction:  nextDirection,
		Source:     "qr",
		OccurredAt: now,
	}
	if err := s.attendanceRepo.Create(log); err != nil {
		return nil, err
	}

	return &dto.ScanQRResponse{
		Accepted:   true,
		Direction:  nextDirection,
		OccurredAt: now.Format(time.RFC3339),
	}, nil
}
