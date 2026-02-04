package services

import (
	"errors"
	"fiber-gorm-app/internal/dto"
	"fiber-gorm-app/internal/models"
	"fiber-gorm-app/internal/repository"
	"time"

	"gorm.io/gorm"
)

type MembershipService interface {
	GetMyMembership(userID uint) (*dto.MembershipResponse, error)
	Renew(userID uint, months int, plan string) (*dto.MembershipResponse, error)
	IsActive(userID uint) (bool, *models.Membership, error)
}

type membershipService struct {
	repo repository.MembershipRepository
}

func NewMembershipService(repo repository.MembershipRepository) MembershipService {
	return &membershipService{repo: repo}
}

func (s *membershipService) GetMyMembership(userID uint) (*dto.MembershipResponse, error) {
	history, err := s.repo.ListByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &dto.MembershipResponse{Current: &dto.MembershipStatusResponse{Status: "expired"}}, nil
		}
		return nil, err
	}

	resp := &dto.MembershipResponse{}
	for i, m := range history {
		item := dto.MembershipHistoryItem{
			Status:   normalizeMembershipStatus(&m),
			Plan:     m.Plan,
			StartsAt: m.StartsAt.Format(time.RFC3339),
			EndsAt:   m.EndsAt.Format(time.RFC3339),
		}
		resp.History = append(resp.History, item)
		if i == 0 {
			resp.Current = &dto.MembershipStatusResponse{
				Status:   item.Status,
				Plan:     item.Plan,
				StartsAt: item.StartsAt,
				EndsAt:   item.EndsAt,
			}
		}
	}

	if resp.Current == nil {
		resp.Current = &dto.MembershipStatusResponse{Status: "expired"}
	}
	return resp, nil
}

func (s *membershipService) Renew(userID uint, months int, plan string) (*dto.MembershipResponse, error) {
	if months <= 0 {
		months = 1
	}
	if plan == "" {
		plan = "standard"
	}

	now := time.Now()
	ends := now.AddDate(0, months, 0)

	// expire any active records
	_ = s.repo.ExpireActiveForUser(userID, now)

	m := &models.Membership{
		UserID:   userID,
		Status:   "active",
		Plan:     plan,
		StartsAt: now,
		EndsAt:   ends,
	}

	if err := s.repo.Create(m); err != nil {
		return nil, err
	}

	return s.GetMyMembership(userID)
}

func (s *membershipService) IsActive(userID uint) (bool, *models.Membership, error) {
	current, err := s.repo.FindLatestByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil, nil
		}
		return false, nil, err
	}

	status := normalizeMembershipStatus(current)
	return status == "active", current, nil
}

func normalizeMembershipStatus(m *models.Membership) string {
	if m == nil {
		return "expired"
	}
	if m.Status != "active" {
		return "expired"
	}
	if time.Now().After(m.EndsAt) {
		return "expired"
	}
	return "active"
}
