package handler

import (
	"encoding/json"
	"fmt"
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
	ch := make(chan ai.ChatCompletionChunk)
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
		case chunk, ok := <-ch:
			if !ok {
				// Channel closed, check for errors
				if err := <-errCh; err != nil {
					// Send error as JSON in data field
					jsonError, _ := json.Marshal(gin.H{"error": err.Error()})
					_, _ = c.Writer.Write([]byte(fmt.Sprintf("data: %s\n\n", jsonError)))
				} else {
					// Send DONE signal
					_, _ = c.Writer.Write([]byte("data: [DONE]\n\n"))
				}
				return false
			}
			// Send chunk as JSON in data field
			jsonChunk, err := json.Marshal(chunk)
			if err != nil {
				// If marshaling fails, send an error
				jsonError, _ := json.Marshal(gin.H{"error": err.Error()})
				_, _ = c.Writer.Write([]byte(fmt.Sprintf("data: %s\n\n", jsonError)))
				return false
			}
			if _, err := c.Writer.Write([]byte(fmt.Sprintf("data: %s\n\n", jsonChunk))); err != nil {
				return false
			}
			c.Writer.Flush()
			return true
		case err := <-errCh:
			if err != nil {
				jsonError, _ := json.Marshal(gin.H{"error": err.Error()})
				_, _ = c.Writer.Write([]byte(fmt.Sprintf("data: %s\n\n", jsonError)))
			}
			return false
		}
	})
}
