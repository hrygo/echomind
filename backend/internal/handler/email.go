package handler

import (
	"net/http"
	"strconv"

	"echomind.com/backend/internal/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type EmailHandler struct {
	db *gorm.DB
}

func NewEmailHandler(db *gorm.DB) *EmailHandler {
	return &EmailHandler{db: db}
}

// ListEmails returns a list of emails.
func (h *EmailHandler) ListEmails(c *gin.Context) {
	var emails []model.Email
	// Simple pagination or limit for MVP
	if err := h.db.Order("date desc").Limit(50).Find(&emails).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, emails)
}

// GetEmail returns a single email by ID.
func (h *EmailHandler) GetEmail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var email model.Email
	if err := h.db.First(&email, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Email not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, email)
}
