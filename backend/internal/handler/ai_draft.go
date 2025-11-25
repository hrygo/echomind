package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hrygo/echomind/internal/service"
)

type AIDraftRequest struct {
	EmailContent string `json:"emailContent" binding:"required"`
	UserPrompt   string `json:"userPrompt" binding:"required"`
}

type AIReplyRequest struct {
	EmailID string  `json:"emailId" binding:"required"`
	Tone    string  `json:"tone,omitempty"`    // "professional", "casual", "friendly", etc.
	Context string  `json:"context,omitempty"`  // "brief", "detailed", etc.
}

type AIReplyResponse struct {
	Reply      string  `json:"reply"`
	Confidence float64 `json:"confidence"`
}

type AIDraftHandler struct {
	aiDraftService *service.AIDraftService
	emailService   *service.EmailService
}

func NewAIDraftHandler(aiDraftService *service.AIDraftService, emailService *service.EmailService) *AIDraftHandler {
	return &AIDraftHandler{
		aiDraftService: aiDraftService,
		emailService:   emailService,
	}
}

// GenerateDraft handles the POST /ai/draft API request.
func (h *AIDraftHandler) GenerateDraft(c *gin.Context) {
	var req AIDraftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	draft, err := h.aiDraftService.GenerateDraftReply(c.Request.Context(), req.EmailContent, req.UserPrompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"draft": draft})
}

// GenerateReply handles the POST /ai/reply API request.
func (h *AIDraftHandler) GenerateReply(c *gin.Context) {
	var req AIReplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("userID").(uuid.UUID)

	// Parse email ID
	emailID, err := uuid.Parse(req.EmailID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email ID"})
		return
	}

	// Get email content
	email, err := h.emailService.GetEmail(c.Request.Context(), userID, emailID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Email not found"})
		return
	}

	// Build user prompt based on tone and context
	userPrompt := h.buildPrompt(req.Tone, req.Context)

	// Generate reply using email body text
	replyText, err := h.aiDraftService.GenerateDraftReply(c.Request.Context(), email.BodyText, userPrompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Calculate confidence (mock for now)
	confidence := 0.85 + (float64(len(replyText)) / 10000.0) // Simple confidence calculation
	if confidence > 0.99 {
		confidence = 0.99
	}

	response := AIReplyResponse{
		Reply:      replyText,
		Confidence: confidence,
	}

	c.JSON(http.StatusOK, response)
}

// buildPrompt creates a user prompt based on tone and context preferences
func (h *AIDraftHandler) buildPrompt(tone, context string) string {
	prompt := "Generate a professional email reply"

	if tone != "" {
		switch tone {
		case "professional":
			prompt = "Generate a professional and formal email reply"
		case "casual":
			prompt = "Generate a casual and friendly email reply"
		case "friendly":
			prompt = "Generate a warm and friendly email reply"
		case "concise":
			prompt = "Generate a concise and brief email reply"
		}
	}

	if context != "" {
		switch context {
		case "brief":
			prompt += " with brief context"
		case "detailed":
			prompt += " with detailed explanation"
		case "urgent":
			prompt += " with urgent and clear action items"
		}
	}

	prompt += ". Include appropriate greeting and closing. Keep it natural and human-like."

	return prompt
}
