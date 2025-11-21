package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hrygo/echomind/internal/middleware"
	"github.com/hrygo/echomind/internal/service"
)

type InsightHandler struct {
	insightService service.InsightService
}

func NewInsightHandler(insightService service.InsightService) *InsightHandler {
	return &InsightHandler{insightService: insightService}
}

// GetNetworkGraph handles the GET /insights/network API request.
func (h *InsightHandler) GetNetworkGraph(c *gin.Context) {
	uuidUserID, ok := middleware.GetUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	graph, err := h.insightService.GetNetworkGraph(c.Request.Context(), uuidUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, graph)
}
