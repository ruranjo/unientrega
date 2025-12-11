package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PasswordReset represents a password reset token
type PasswordReset struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Token     string    `gorm:"uniqueIndex;not null;size:255" json:"token"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	Used      bool      `gorm:"default:false" json:"used"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relationships
	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
}

// TableName specifies the table name for PasswordReset model
func (PasswordReset) TableName() string {
	return "password_resets"
}

// BeforeCreate is a GORM hook that runs before creating a password reset
func (p *PasswordReset) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

// IsExpired checks if the reset token has expired
func (p *PasswordReset) IsExpired() bool {
	return time.Now().After(p.ExpiresAt)
}

// IsValid checks if the token is valid (not expired and not used)
func (p *PasswordReset) IsValid() bool {
	return !p.Used && !p.IsExpired()
}
