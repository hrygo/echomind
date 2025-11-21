package handler

import (
	"errors"
	"net/http"

	"github.com/hrygo/echomind/internal/middleware"
	"github.com/hrygo/echomind/internal/service"
	"github.com/gin-gonic/gin"
)

// SyncHandler handles email synchronization requests.
type SyncHandler struct {
	syncService *service.SyncService
}

// NewSyncHandler creates a new SyncHandler.
func NewSyncHandler(syncService *service.SyncService) *SyncHandler {
	return &SyncHandler{
		syncService: syncService,
	}
}

// SyncEmails handles the email synchronization request for the authenticated user.
func (h *SyncHandler) SyncEmails(c *gin.Context) {
	userID, ok := middleware.GetUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	err := h.syncService.SyncEmails(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, service.ErrAccountNotConfigured) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Please configure your email account in Settings first."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sync initiated successfully"})
}
