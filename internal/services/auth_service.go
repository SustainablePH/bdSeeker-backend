package services

import (
	"errors"
	"time"

	"github.com/bishworup11/bdSeeker-backend/internal/config"
	"github.com/bishworup11/bdSeeker-backend/internal/models"
	"github.com/bishworup11/bdSeeker-backend/internal/repositories"
	"github.com/bishworup11/bdSeeker-backend/pkg/utils"
	"gorm.io/gorm"
)

type AuthService struct {
	userRepo *repositories.UserRepository
}

func NewAuthService(userRepo *repositories.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	FullName string `json:"full_name" validate:"required"`
	Role     string `json:"role" validate:"required,oneof=developer company admin"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
}

func (s *AuthService) Register(req *RegisterRequest) (*AuthResponse, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.FindByEmail(req.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		Email:        req.Email,
		PasswordHash: hashedPassword,
		FullName:     req.FullName,
		Role:         req.Role,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// Generate token
	expiry, _ := time.ParseDuration(config.AppConfig.JWTExpiry)
	token, err := utils.GenerateToken(user.ID, user.Email, user.Role, config.AppConfig.JWTSecret, expiry)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

func (s *AuthService) Login(req *LoginRequest) (*AuthResponse, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid email or password")
		}
		return nil, err
	}

	// Verify password
	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		return nil, errors.New("invalid email or password")
	}

	// Generate token
	expiry, _ := time.ParseDuration(config.AppConfig.JWTExpiry)
	token, err := utils.GenerateToken(user.ID, user.Email, user.Role, config.AppConfig.JWTSecret, expiry)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

func (s *AuthService) GetUserByID(userID uint) (*models.User, error) {
	return s.userRepo.FindByID(userID)
}
