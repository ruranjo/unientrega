package repository

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/ruranjo/unientrega/internal/models"
)

// PasswordResetRepository handles database operations for password resets
type PasswordResetRepository struct {
	db *gorm.DB
}

// NewPasswordResetRepository creates a new password reset repository
func NewPasswordResetRepository(db *gorm.DB) *PasswordResetRepository {
	return &PasswordResetRepository{db: db}
}

// Create creates a new password reset token
func (r *PasswordResetRepository) Create(reset *models.PasswordReset) error {
	return r.db.Create(reset).Error
}

// GetByToken finds a password reset by token
func (r *PasswordResetRepository) GetByToken(token string) (*models.PasswordReset, error) {
	var reset models.PasswordReset
	err := r.db.Where("token = ?", token).First(&reset).Error
	if err != nil {
		return nil, err
	}
	return &reset, nil
}

// GetByUserID finds all password resets for a user
func (r *PasswordResetRepository) GetByUserID(userID uuid.UUID) ([]*models.PasswordReset, error) {
	var resets []*models.PasswordReset
	err := r.db.Where("user_id = ?", userID).Find(&resets).Error
	return resets, err
}

// MarkAsUsed marks a password reset token as used
func (r *PasswordResetRepository) MarkAsUsed(token string) error {
	return r.db.Model(&models.PasswordReset{}).
		Where("token = ?", token).
		Update("used", true).Error
}

// DeleteExpired deletes all expired password reset tokens
func (r *PasswordResetRepository) DeleteExpired() error {
	return r.db.Where("expires_at < ?", time.Now()).Delete(&models.PasswordReset{}).Error
}

// DeleteByUserID deletes all password reset tokens for a user
func (r *PasswordResetRepository) DeleteByUserID(userID uuid.UUID) error {
	return r.db.Where("user_id = ?", userID).Delete(&models.PasswordReset{}).Error
}
