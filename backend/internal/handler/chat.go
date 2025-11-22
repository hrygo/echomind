package handler

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hrygo/echomind/internal/service"
	"github.com/hrygo/echomind/pkg/ai"
)

type ChatHandler struct {
	chatService *service.ChatService
}

func NewChatHandler(chatService *service.ChatService) *ChatHandler {
	return &ChatHandler{
		chatService: chatService,
	}
}

type ChatRequest struct {
	Messages []ai.Message `json:"messages" binding:"required"`
}

func (h *ChatHandler) StreamChat(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get User ID from context (set by AuthMiddleware)
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id"})
		return
	}

	// Set headers for SSE
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")

	// Create a channel to receive stream
	ch := make(chan string)
	errCh := make(chan error)

	go func() {
		err := h.chatService.StreamChat(c.Request.Context(), userID, req.Messages, ch)
		if err != nil {
			errCh <- err
		}
		close(errCh)
	}()

	c.Stream(func(w io.Writer) bool {
		select {
		case msg, ok := <-ch:
			if !ok {
				// Channel closed, check for errors
				if err := <-errCh; err != nil {
					c.SSEvent("error", err.Error())
				}
				return false
			}
			c.SSEvent("message", msg)
			return true
		case err := <-errCh:
			if err != nil {
				c.SSEvent("error", err.Error())
			}
			return false
		}
	})
}
