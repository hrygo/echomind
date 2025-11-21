package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hrygo/echomind/internal/service"
)

type AIDraftRequest struct {
	EmailContent string `json:"emailContent" binding:"required"`
	UserPrompt   string `json:"userPrompt" binding:"required"`
}

type AIDraftHandler struct {
	aiDraftService *service.AIDraftService
}

func NewAIDraftHandler(aiDraftService *service.AIDraftService) *AIDraftHandler {
	return &AIDraftHandler{aiDraftService: aiDraftService}
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
