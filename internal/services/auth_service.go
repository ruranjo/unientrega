package services

import (
	"errors"

	"github.com/google/uuid"
	"github.com/ruranjo/unientrega/internal/models"
	"github.com/ruranjo/unientrega/internal/utils"
)

// AuthService handles authentication business logic
type AuthService struct {
	userService *UserService
}

// NewAuthService creates a new auth service
func NewAuthService(userService *UserService) *AuthService {
	return &AuthService{
		userService: userService,
	}
}

// RegisterRequest represents registration data
type RegisterRequest struct {
	Email     string      `json:"email" binding:"required,email"`
	Password  string      `json:"password" binding:"required,min=6"`
	FirstName string      `json:"first_name"`
	LastName  string      `json:"last_name"`
	Role      models.Role `json:"role"`
}

// LoginRequest represents login credentials
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse represents authentication response with tokens
type AuthResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         *models.User `json:"user"`
}

// Register creates a new user account
func (s *AuthService) Register(req *RegisterRequest) (*AuthResponse, error) {
	// Create user
	user := &models.User{
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      req.Role,
		IsActive:  true,
	}

	// Default to client role if not specified or invalid
	if user.Role == "" || !user.Role.IsValid() {
		user.Role = models.RoleClient
	}

	err := s.userService.CreateUser(user, req.Password)
	if err != nil {
		return nil, err
	}

	// Generate tokens
	accessToken, err := utils.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

// Login authenticates a user and returns tokens
func (s *AuthService) Login(req *LoginRequest) (*AuthResponse, error) {
	// Authenticate user
	user, err := s.userService.AuthenticateUser(req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	// Generate tokens
	accessToken, err := utils.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

// RefreshToken generates a new access token from a refresh token
func (s *AuthService) RefreshToken(refreshToken string) (string, error) {
	// Validate refresh token
	userID, err := utils.ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", errors.New("invalid refresh token")
	}

	// Get user
	user, err := s.userService.GetUserByID(userID)
	if err != nil {
		return "", errors.New("user not found")
	}

	if !user.IsActive {
		return "", errors.New("user account is inactive")
	}

	// Generate new access token
	accessToken, err := utils.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

// GetUserByID retrieves a user by ID
func (s *AuthService) GetUserByID(userID uuid.UUID) (*models.User, error) {
	return s.userService.GetUserByID(userID)
}
