package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hrygo/echomind/internal/middleware"
	"github.com/hrygo/echomind/internal/model"
	"gorm.io/gorm"
)

// EmailServicer defines the interface for the email service that the handler depends on.
type EmailServicer interface {
	ListEmails(ctx context.Context, userID uuid.UUID, limit, offset int, contextID, folder, category, filter string) ([]model.Email, error)
	GetEmail(ctx context.Context, userID, emailID uuid.UUID) (*model.Email, error)
	DeleteAllUserEmails(ctx context.Context, userID uuid.UUID) error
}

type EmailHandler struct {
	emailService EmailServicer
}

func NewEmailHandler(emailService EmailServicer) *EmailHandler {
	return &EmailHandler{emailService: emailService}
}

// ListEmails returns a list of emails for the authenticated user.
func (h *EmailHandler) ListEmails(c *gin.Context) {
	userID, ok := middleware.GetUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	// Pagination and filter parameters
	limitStr := c.DefaultQuery("limit", "50")
	offsetStr := c.DefaultQuery("offset", "0")
	contextID := c.Query("context_id")
	folder := c.Query("folder")
	category := c.Query("category")
	filter := c.Query("filter")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset parameter"})
		return
	}

	emails, err := h.emailService.ListEmails(c.Request.Context(), userID, limit, offset, contextID, folder, category, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, emails)
}

// GetEmail returns a single email by ID for the authenticated user.
func (h *EmailHandler) GetEmail(c *gin.Context) {
	userID, ok := middleware.GetUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	idStr := c.Param("id")
	emailID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email ID format"})
		return
	}

	email, err := h.emailService.GetEmail(c.Request.Context(), userID, emailID)
	if err != nil {
		if err == gorm.ErrRecordNotFound || email == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Email not found or not accessible"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, email)
}

// DeleteAllEmails deletes all emails for the authenticated user.
func (h *EmailHandler) DeleteAllEmails(c *gin.Context) {
	userID, ok := middleware.GetUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	if err := h.emailService.DeleteAllUserEmails(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete all emails", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "All emails deleted successfully"})
}
