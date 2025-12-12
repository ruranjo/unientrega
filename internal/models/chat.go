package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ChatMessage represents a chat message between users
type ChatMessage struct {
	ID         uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	OrderID    uuid.UUID      `gorm:"type:uuid;not null" json:"order_id"`
	SenderID   uuid.UUID      `gorm:"type:uuid;not null" json:"sender_id"`
	ReceiverID uuid.UUID      `gorm:"type:uuid;not null" json:"receiver_id"`
	Content    string         `gorm:"type:text;not null" json:"content"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for ChatMessage model
func (ChatMessage) TableName() string {
	return "chat_messages"
}

// BeforeCreate is a GORM hook that runs before creating a message
func (m *ChatMessage) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}
