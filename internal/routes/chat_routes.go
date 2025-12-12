package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ruranjo/unientrega/internal/handlers"
)

func SetupChatRoutes(rg *gin.RouterGroup, handler *handlers.ChatHandler) {
	chat := rg.Group("/chat")
	{
		chat.GET("/ws", handler.HandleWebSocket)
		chat.GET("/history/:orderID", handler.GetHistory)
	}
}
