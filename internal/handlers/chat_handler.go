package handlers

import (
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/ruranjo/unientrega/internal/services"
)

type ChatHandler struct {
	chatService *services.ChatService
	clients     map[uuid.UUID]*websocket.Conn
	mu          sync.Mutex
	upgrader    websocket.Upgrader
}

func NewChatHandler(chatService *services.ChatService) *ChatHandler {
	return &ChatHandler{
		chatService: chatService,
		clients:     make(map[uuid.UUID]*websocket.Conn),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins for now
			},
		},
	}
}

func (h *ChatHandler) HandleWebSocket(c *gin.Context) {
	userIDStr := c.Query("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	h.mu.Lock()
	h.clients[userID] = conn
	h.mu.Unlock()

	defer func() {
		h.mu.Lock()
		delete(h.clients, userID)
		h.mu.Unlock()
		conn.Close()
	}()

	for {
		var msg struct {
			OrderID    uuid.UUID `json:"order_id"`
			ReceiverID uuid.UUID `json:"receiver_id"`
			Content    string    `json:"content"`
		}

		if err := conn.ReadJSON(&msg); err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}

		// Save message to database
		savedMsg, err := h.chatService.SendMessage(msg.OrderID, userID, msg.ReceiverID, msg.Content)
		if err != nil {
			log.Printf("Error saving message: %v", err)
			continue
		}

		// Send to receiver if connected
		h.mu.Lock()
		receiverConn, ok := h.clients[msg.ReceiverID]
		h.mu.Unlock()

		if ok {
			if err := receiverConn.WriteJSON(savedMsg); err != nil {
				log.Printf("Error sending message to receiver: %v", err)
			}
		}
	}
}

func (h *ChatHandler) GetHistory(c *gin.Context) {
	orderIDStr := c.Param("orderID")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	messages, err := h.chatService.GetChatHistory(orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch chat history"})
		return
	}

	c.JSON(http.StatusOK, messages)
}
