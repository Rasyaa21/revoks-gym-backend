package services

import (
	"errors"
	"fiber-gorm-app/internal/dto"
	"fiber-gorm-app/internal/repository"
	"fiber-gorm-app/internal/utils"
	"os"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService interface {
	Login(req *dto.LoginRequest) (*dto.LoginResponse, error)
}

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

func (s *authService) Login(req *dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	ttl := accessTokenTTL()
	token, err := utils.NewAccessToken(user.ID, ttl)
	if err != nil {
		return nil, err
	}

	expiresIn := int64(ttl.Seconds())
	return &dto.LoginResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   expiresIn,
		User: &dto.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			Phone:     user.Phone,
			Address:   user.Address,
			IsActive:  user.IsActive,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
			UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
		},
	}, nil
}

func accessTokenTTL() time.Duration {
	minutesStr := os.Getenv("JWT_ACCESS_TTL_MINUTES")
	if minutesStr == "" {
		return 24 * time.Hour
	}
	minutes, err := strconv.Atoi(minutesStr)
	if err != nil || minutes <= 0 {
		return 24 * time.Hour
	}
	return time.Duration(minutes) * time.Minute
}
