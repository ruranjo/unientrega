package repository

import (
	"github.com/google/uuid"
	"github.com/ruranjo/unientrega/internal/models"
	"gorm.io/gorm"
)

type ChatRepository struct {
	db *gorm.DB
}

func NewChatRepository(db *gorm.DB) *ChatRepository {
	return &ChatRepository{db: db}
}

func (r *ChatRepository) CreateMessage(message *models.ChatMessage) error {
	return r.db.Create(message).Error
}

func (r *ChatRepository) GetMessagesByOrder(orderID uuid.UUID) ([]models.ChatMessage, error) {
	var messages []models.ChatMessage
	err := r.db.Where("order_id = ?", orderID).Order("created_at asc").Find(&messages).Error
	return messages, err
}
