package services

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/ruranjo/unientrega/internal/models"
	"github.com/ruranjo/unientrega/internal/repository"
)

// UserService handles business logic for users
type UserService struct {
	userRepo          *repository.UserRepository
	passwordResetRepo *repository.PasswordResetRepository
}

// NewUserService creates a new user service
func NewUserService(userRepo *repository.UserRepository, passwordResetRepo *repository.PasswordResetRepository) *UserService {
	return &UserService{
		userRepo:          userRepo,
		passwordResetRepo: passwordResetRepo,
	}
}

// CreateUser creates a new user with hashed password
func (s *UserService) CreateUser(user *models.User, plainPassword string) error {
	// Check if email already exists
	exists, err := s.userRepo.ExistsByEmail(user.Email)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("email already exists")
	}

	// Hash password
	hashedPassword, err := HashPassword(plainPassword)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	// Validate role
	if !user.Role.IsValid() {
		user.Role = models.RoleClient // Default to client
	}

	return s.userRepo.Create(user)
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(id uuid.UUID) (*models.User, error) {
	return s.userRepo.GetByID(id)
}

// GetUserByEmail retrieves a user by email
func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	return s.userRepo.GetByEmail(email)
}

// UpdateUser updates a user
func (s *UserService) UpdateUser(user *models.User) error {
	// Validate role if changed
	if !user.Role.IsValid() {
		return errors.New("invalid role")
	}
	return s.userRepo.Update(user)
}

// UpdatePassword updates a user's password
func (s *UserService) UpdatePassword(userID uuid.UUID, newPassword string) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	hashedPassword, err := HashPassword(newPassword)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	return s.userRepo.Update(user)
}

// DeleteUser soft deletes a user
func (s *UserService) DeleteUser(id uuid.UUID) error {
	return s.userRepo.Delete(id)
}

// ListUsers returns a list of users with pagination
func (s *UserService) ListUsers(limit, offset int, role models.Role) ([]*models.User, error) {
	return s.userRepo.List(limit, offset, role)
}

// AuthenticateUser validates user credentials
func (s *UserService) AuthenticateUser(email, password string) (*models.User, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !user.IsActive {
		return nil, errors.New("user account is inactive")
	}

	if !CheckPassword(user.Password, password) {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

// GetFullName returns the user's full name
func (s *UserService) GetFullName(user *models.User) string {
	if user.FirstName == "" && user.LastName == "" {
		return user.Email
	}
	return user.FirstName + " " + user.LastName
}

// HasPermission checks if a user has a specific permission
func (s *UserService) HasPermission(user *models.User, permission string) bool {
	switch permission {
	case "manage_users":
		return user.Role == models.RoleSuperUser
	case "create_delivery":
		return user.Role == models.RoleSuperUser || user.Role == models.RoleStore
	case "update_delivery_status":
		return user.Role == models.RoleSuperUser || user.Role == models.RoleDelivery
	case "place_order":
		return user.Role == models.RoleSuperUser || user.Role == models.RoleClient
	case "view_all_deliveries":
		return user.Role == models.RoleSuperUser || user.Role == models.RoleStore
	case "view_own_deliveries":
		return true // All roles can view their own deliveries
	default:
		return false
	}
}

// IsSuperUser checks if the user is a superuser
func (s *UserService) IsSuperUser(user *models.User) bool {
	return user.Role == models.RoleSuperUser
}

// Password utility functions

// HashPassword hashes a plain text password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword compares a plain text password with the hashed password
func CheckPassword(hashedPassword, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}

// Password Reset Methods

// GeneratePasswordResetToken generates a password reset token for a user
func (s *UserService) GeneratePasswordResetToken(email string) (*models.PasswordReset, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Delete any existing unused tokens for this user
	s.passwordResetRepo.DeleteByUserID(user.ID)

	// Generate random token
	token := uuid.New().String()

	// Create password reset record
	reset := &models.PasswordReset{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(1 * time.Hour),
		Used:      false,
	}

	err = s.passwordResetRepo.Create(reset)
	if err != nil {
		return nil, err
	}

	return reset, nil
}

// ValidateResetToken checks if a reset token is valid
func (s *UserService) ValidateResetToken(token string) (*models.PasswordReset, error) {
	reset, err := s.passwordResetRepo.GetByToken(token)
	if err != nil {
		return nil, errors.New("invalid token")
	}

	if !reset.IsValid() {
		return nil, errors.New("token has expired or been used")
	}

	return reset, nil
}

// ResetPasswordWithToken resets a user's password using a valid token
func (s *UserService) ResetPasswordWithToken(token, newPassword string) error {
	reset, err := s.ValidateResetToken(token)
	if err != nil {
		return err
	}

	// Update user password
	err = s.UpdatePassword(reset.UserID, newPassword)
	if err != nil {
		return err
	}

	// Mark token as used
	return s.passwordResetRepo.MarkAsUsed(token)
}
