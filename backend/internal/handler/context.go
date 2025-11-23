package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hrygo/echomind/internal/model"
	"github.com/hrygo/echomind/internal/service"
)

type ContextHandler struct {
	contextService *service.ContextService
}

func NewContextHandler(contextService *service.ContextService) *ContextHandler {
	return &ContextHandler{contextService: contextService}
}

// CreateContext creates a new context.
func (h *ContextHandler) CreateContext(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	var input model.ContextInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, err := h.contextService.CreateContext(userID, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, ctx)
}

// ListContexts returns all contexts for a user.
func (h *ContextHandler) ListContexts(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	contexts, err := h.contextService.ListContexts(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, contexts)
}

// UpdateContext updates a context.
func (h *ContextHandler) UpdateContext(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid context ID"})
		return
	}

	var input model.ContextInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, err := h.contextService.UpdateContext(id, userID, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ctx)
}

// DeleteContext deletes a context.
func (h *ContextHandler) DeleteContext(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid context ID"})
		return
	}

	if err := h.contextService.DeleteContext(id, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
