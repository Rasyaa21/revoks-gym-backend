package services

import (
	"errors"
	"fiber-gorm-app/internal/dto"
	"fiber-gorm-app/internal/models"
	"fiber-gorm-app/internal/repository"
	"math"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService interface {
	CreateUser(req *dto.CreateUserRequest) (*dto.UserResponse, error)
	GetAllUsers(page, perPage int) ([]dto.UserResponse, *dto.PaginationResponse, error)
	GetUserByID(id uuid.UUID) (*dto.UserResponse, error)
	UpdateUser(id uuid.UUID, req *dto.UpdateUserRequest) (*dto.UserResponse, error)
	DeleteUser(id uuid.UUID) error
	ChangePassword(id uuid.UUID, req *dto.ChangePasswordRequest) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(req *dto.CreateUserRequest) (*dto.UserResponse, error) {
	// Check if email already exists (if provided)
	if req.Email != "" {
		existingUser, _ := s.repo.FindByEmail(req.Email)
		if existingUser != nil {
			return nil, errors.New("email already registered")
		}
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Role:         models.UserRole(req.Role),
		FullName:     req.FullName,
		PasswordHash: string(hashedPassword),
		Status:       models.UserStatusActive,
	}

	// Set optional fields if provided
	if req.Email != "" {
		user.Email = &req.Email
	}
	if req.Phone != "" {
		user.Phone = &req.Phone
	}
	if req.PhotoURL != "" {
		user.PhotoURL = &req.PhotoURL
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return s.toUserResponse(user), nil
}

func (s *userService) GetAllUsers(page, perPage int) ([]dto.UserResponse, *dto.PaginationResponse, error) {
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 10
	}

	users, total, err := s.repo.FindAll(page, perPage)
	if err != nil {
		return nil, nil, err
	}

	var responses []dto.UserResponse
	for _, user := range users {
		responses = append(responses, *s.toUserResponse(&user))
	}

	pagination := &dto.PaginationResponse{
		CurrentPage: page,
		PerPage:     perPage,
		TotalPages:  int(math.Ceil(float64(total) / float64(perPage))),
		TotalItems:  total,
	}

	return responses, pagination, nil
}

func (s *userService) GetUserByID(id uuid.UUID) (*dto.UserResponse, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return s.toUserResponse(user), nil
}

func (s *userService) UpdateUser(id uuid.UUID, req *dto.UpdateUserRequest) (*dto.UserResponse, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	// Update fields if provided
	if req.FullName != "" {
		user.FullName = req.FullName
	}
	if req.Email != "" {
		// Check if new email is already used by another user
		existingUser, _ := s.repo.FindByEmail(req.Email)
		if existingUser != nil && existingUser.UserID != id {
			return nil, errors.New("email already used by another user")
		}
		user.Email = &req.Email
	}
	if req.Phone != "" {
		user.Phone = &req.Phone
	}
	if req.PhotoURL != "" {
		user.PhotoURL = &req.PhotoURL
	}
	if req.Status != "" {
		user.Status = models.UserStatus(req.Status)
	}

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	return s.toUserResponse(user), nil
}

func (s *userService) DeleteUser(id uuid.UUID) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	return s.repo.Delete(id)
}

func (s *userService) ChangePassword(id uuid.UUID, req *dto.ChangePasswordRequest) error {
	user, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.OldPassword)); err != nil {
		return errors.New("old password is incorrect")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.PasswordHash = string(hashedPassword)
	return s.repo.Update(user)
}

func (s *userService) toUserResponse(user *models.User) *dto.UserResponse {
	resp := &dto.UserResponse{
		UserID:    user.UserID.String(),
		Role:      string(user.Role),
		FullName:  user.FullName,
		Status:    string(user.Status),
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	}
	
	if user.Email != nil {
		resp.Email = *user.Email
	}
	if user.Phone != nil {
		resp.Phone = *user.Phone
	}
	if user.PhotoURL != nil {
		resp.PhotoURL = *user.PhotoURL
	}
	
	return resp
}
