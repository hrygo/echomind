package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hrygo/echomind/internal/middleware"
	"github.com/hrygo/echomind/internal/service"
	"go.uber.org/zap"
)

type Searcher interface {
	Search(ctx context.Context, userID uuid.UUID, query string, filters service.SearchFilters, limit int) ([]service.SearchResult, error)
}

type SearchHandler struct {
	searchService Searcher
	logger        *zap.SugaredLogger
}

func NewSearchHandler(searchService Searcher, logger *zap.SugaredLogger) *SearchHandler {
	return &SearchHandler{
		searchService: searchService,
		logger:        logger,
	}
}

func (h *SearchHandler) Search(c *gin.Context) {
	start := time.Now()

	// Extract user_id from context (set by auth middleware)
	userIDValue, exists := c.Get(middleware.ContextUserIDKey)
	if !exists {
		h.logger.Warnw("Search attempt without authentication")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	userID, ok := userIDValue.(uuid.UUID)
	if !ok {
		h.logger.Errorw("Invalid user ID format in context", "userIDValue", userIDValue)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID format"})
		return
	}

	// Get query parameter
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter 'q' is required"})
		return
	}

	// Parse filters
	var filters service.SearchFilters
	filters.Sender = c.Query("sender")

	if startDateStr := c.Query("start_date"); startDateStr != "" {
		if t, err := time.Parse(time.DateOnly, startDateStr); err == nil {
			filters.StartDate = &t
		} else {
			h.logger.Warnw("Invalid start_date format", "start_date", startDateStr)
		}
	}

	if endDateStr := c.Query("end_date"); endDateStr != "" {
		if t, err := time.Parse(time.DateOnly, endDateStr); err == nil {
			// Add 23:59:59 to include the entire end date
			t = t.Add(24*time.Hour - time.Nanosecond)
			filters.EndDate = &t
		} else {
			h.logger.Warnw("Invalid end_date format", "end_date", endDateStr)
		}
	}

	// Get limit parameter (optional, default to 10)
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 100 {
		limit = 10
	}

	h.logger.Infow("Search request", "userID", userID, "query", query, "filters", filters, "limit", limit)

	// Perform search
	results, err := h.searchService.Search(c.Request.Context(), userID, query, filters, limit)
	duration := time.Since(start)

	if err != nil {
		h.logger.Errorw("Search failed", "userID", userID, "query", query, "error", err, "duration", duration)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "search failed", "details": err.Error()})
		return
	}

	h.logger.Infow("Search completed", "userID", userID, "query", query, "results", len(results), "duration", duration)

	c.JSON(http.StatusOK, gin.H{
		"query":   query,
		"results": results,
		"count":   len(results),
	})
}
