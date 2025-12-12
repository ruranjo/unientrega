package services

import (
	"github.com/google/uuid"
	"github.com/ruranjo/unientrega/internal/models"
	"github.com/ruranjo/unientrega/internal/repository"
)

type ChatService struct {
	chatRepo *repository.ChatRepository
}

func NewChatService(chatRepo *repository.ChatRepository) *ChatService {
	return &ChatService{chatRepo: chatRepo}
}

func (s *ChatService) SendMessage(orderID, senderID, receiverID uuid.UUID, content string) (*models.ChatMessage, error) {
	message := &models.ChatMessage{
		OrderID:    orderID,
		SenderID:   senderID,
		ReceiverID: receiverID,
		Content:    content,
	}

	if err := s.chatRepo.CreateMessage(message); err != nil {
		return nil, err
	}

	return message, nil
}

func (s *ChatService) GetChatHistory(orderID uuid.UUID) ([]models.ChatMessage, error) {
	return s.chatRepo.GetMessagesByOrder(orderID)
}
