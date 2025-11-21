package handler

import (
	"net/http"
	"strconv"

	"github.com/hrygo/echomind/internal/middleware"
	"github.com/hrygo/echomind/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EmailHandler struct {
	emailService *service.EmailService
}

func NewEmailHandler(emailService *service.EmailService) *EmailHandler {
	return &EmailHandler{emailService: emailService}
}

// ListEmails returns a list of emails for the authenticated user.
func (h *EmailHandler) ListEmails(c *gin.Context) {
	userID, ok := middleware.GetUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	// Pagination parameters
	limitStr := c.DefaultQuery("limit", "50")
	offsetStr := c.DefaultQuery("offset", "0")

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

	emails, err := h.emailService.ListEmails(c.Request.Context(), userID, limit, offset)
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
