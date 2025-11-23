package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hrygo/echomind/internal/service"
)

type ActionHandler struct {
	actionService *service.ActionService
}

func NewActionHandler(actionService *service.ActionService) *ActionHandler {
	return &ActionHandler{actionService: actionService}
}

type ApproveRequest struct {
	EmailID string `json:"email_id" binding:"required"`
}

type SnoozeRequest struct {
	EmailID  string `json:"email_id" binding:"required"`
	Duration string `json:"duration"` // "4h", "tomorrow", or ISO timestamp
}

type DismissRequest struct {
	EmailID string `json:"email_id" binding:"required"`
}

// ApproveEmail handles the approval action (archive/done).
func (h *ActionHandler) ApproveEmail(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	var req ApproveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	emailID, err := uuid.Parse(req.EmailID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email id"})
		return
	}

	if err := h.actionService.ApproveEmail(c.Request.Context(), userID, emailID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "approved"})
}

// SnoozeEmail handles snoozing an email.
func (h *ActionHandler) SnoozeEmail(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	var req SnoozeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	emailID, err := uuid.Parse(req.EmailID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email id"})
		return
	}

	// Calculate Snooze Time
	var snoozeUntil time.Time
	now := time.Now()

	switch req.Duration {
	case "4h":
		snoozeUntil = now.Add(4 * time.Hour)
	case "tomorrow":
		// Tomorrow at 9 AM
		snoozeUntil = now.Add(24 * time.Hour)
		snoozeUntil = time.Date(snoozeUntil.Year(), snoozeUntil.Month(), snoozeUntil.Day(), 9, 0, 0, 0, snoozeUntil.Location())
	case "next_week":
		snoozeUntil = now.Add(7 * 24 * time.Hour)
		snoozeUntil = time.Date(snoozeUntil.Year(), snoozeUntil.Month(), snoozeUntil.Day(), 9, 0, 0, 0, snoozeUntil.Location())
	default:
		// Try parsing as specific time if needed, or default to 4h
		parsed, err := time.Parse(time.RFC3339, req.Duration)
		if err == nil {
			snoozeUntil = parsed
		} else {
			snoozeUntil = now.Add(4 * time.Hour)
		}
	}

	if err := h.actionService.SnoozeEmail(c.Request.Context(), userID, emailID, snoozeUntil); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "snoozed", "until": snoozeUntil})
}

// DismissEmail removes the email from the smart feed.
func (h *ActionHandler) DismissEmail(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	var req DismissRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	emailID, err := uuid.Parse(req.EmailID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email id"})
		return
	}

	if err := h.actionService.DismissEmail(c.Request.Context(), userID, emailID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "dismissed"})
}
