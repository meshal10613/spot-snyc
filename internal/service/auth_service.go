package service

import (
	"errors"
	"fmt"
	"spot-sync/internal/config"
	"spot-sync/internal/dto"
	"spot-sync/internal/models"
	"spot-sync/internal/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// AuthService defines the auth business logic contract.
type AuthService interface {
	Register(req *dto.RegisterRequest) (*dto.UserResponse, error)
	Login(req *dto.LoginRequest) (*dto.LoginResponse, error)
}

type authService struct {
	repo      repository.AuthRepository
	jwtSecret string
}

// NewAuthService creates a new auth service with injected dependencies.
func NewAuthService(repo repository.AuthRepository, cfg *config.Config) AuthService {
	return &authService{
		repo:      repo,
		jwtSecret: cfg.JWTSecret,
	}
}

func (s *authService) Register(req *dto.RegisterRequest) (*dto.UserResponse, error) {
	// Check if email is already registered
	existing, err := s.repo.FindUserByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if existing != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Default role to "driver" if not provided
	role := models.RoleDriver
	if req.Role != "" {
		role = models.Role(req.Role)
	}

	user := &models.User{
		Name:  req.Name,
		Email: req.Email,
		Role:  role,
	}

	// Hash password with bcrypt (cost 12)
	if err := user.HashPassword(req.Password); err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	if err := s.repo.CreateUser(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      string(user.Role),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (s *authService) Login(req *dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.repo.FindUserByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if user == nil {
		return nil, errors.New("invalid email or password")
	}

	// Compare bcrypt hash
	if err := user.CheckPassword(req.Password); err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Generate JWT with user_id and role in the payload
	token, err := s.generateToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &dto.LoginResponse{
		Token: token,
		User: dto.LoginUserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  string(user.Role),
		},
	}, nil
}

// generateToken creates a signed JWT containing user_id and role.
func (s *authService) generateToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role":    string(user.Role),
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}
