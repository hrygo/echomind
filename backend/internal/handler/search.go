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
	"github.com/hrygo/echomind/pkg/logger"
)

type Searcher interface {
	Search(ctx context.Context, userID uuid.UUID, query string, filters service.SearchFilters, limit int) ([]service.SearchResult, error)
}

type SearchClusterer interface {
	ClusterResults(results []service.SearchResult, clusterType service.ClusterType) []service.SearchCluster
}

type SearchSummarizer interface {
	GenerateSummary(ctx context.Context, results []service.SearchResult, query string) (*service.SearchResultsSummary, error)
	GenerateQuickSummary(results []service.SearchResult) *service.SearchResultsSummary
}

type SearchHandler struct {
	searchService     Searcher
	clusteringService SearchClusterer
	summaryService    SearchSummarizer
	logger            logger.Logger
}

func NewSearchHandler(searchService Searcher, clusteringService SearchClusterer, summaryService SearchSummarizer, log logger.Logger) *SearchHandler {
	return &SearchHandler{
		searchService:     searchService,
		clusteringService: clusteringService,
		summaryService:    summaryService,
		logger:            log,
	}
}

func (h *SearchHandler) Search(c *gin.Context) {
	start := time.Now()

	// Extract user_id from context (set by auth middleware)
	userIDValue, exists := c.Get(middleware.ContextUserIDKey)
	if !exists {
		h.logger.Warn("Search attempt without authentication")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	userID, ok := userIDValue.(uuid.UUID)
	if !ok {
		h.logger.Error("Invalid user ID format in context", logger.Any("userIDValue", userIDValue))
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
	if contextIDStr := c.Query("context_id"); contextIDStr != "" {
		if contextID, err := uuid.Parse(contextIDStr); err == nil {
			filters.ContextID = &contextID
		} else {
			h.logger.Warn("Invalid context_id format", logger.String("context_id", contextIDStr))
		}
	}

	if startDateStr := c.Query("start_date"); startDateStr != "" {
		if t, err := time.Parse(time.DateOnly, startDateStr); err == nil {
			filters.StartDate = &t
		} else {
			h.logger.Warn("Invalid start_date format", logger.String("start_date", startDateStr))
		}
	}

	if endDateStr := c.Query("end_date"); endDateStr != "" {
		if t, err := time.Parse(time.DateOnly, endDateStr); err == nil {
			// Add 23:59:59 to include the entire end date
			t = t.Add(24*time.Hour - time.Nanosecond)
			filters.EndDate = &t
		} else {
			h.logger.Warn("Invalid end_date format", logger.String("end_date", endDateStr))
		}
	}

	// Get limit parameter (optional, default to 10)
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 100 {
		limit = 10
	}

	// Get enhancement flags
	enableClustering := c.DefaultQuery("enable_clustering", "false") == "true"
	enableSummary := c.DefaultQuery("enable_summary", "false") == "true"
	clusterTypeStr := c.DefaultQuery("cluster_type", "sender") // sender, time, topic

	h.logger.Info("Search request",
		logger.Any("userID", userID),
		logger.String("query", query),
		logger.Any("filters", filters),
		logger.Int("limit", limit),
	)

	// Perform search
	results, err := h.searchService.Search(c.Request.Context(), userID, query, filters, limit)
	duration := time.Since(start)

	if err != nil {
		h.logger.Error("Search failed",
			logger.Any("userID", userID),
			logger.String("query", query),
			logger.Error(err),
			logger.Duration("duration", duration),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "search failed", "details": err.Error()})
		return
	}

	h.logger.Info("Search completed",
		logger.Any("userID", userID),
		logger.String("query", query),
		logger.Int("results", len(results)),
		logger.Duration("duration", duration),
	)

	// Prepare response with base data
	response := gin.H{
		"query":   query,
		"results": results,
		"count":   len(results),
	}

	// Add clustering if enabled
	if enableClustering && h.clusteringService != nil && len(results) > 0 {
		var clusterType service.ClusterType
		switch clusterTypeStr {
		case "time":
			clusterType = service.ClusterByTime
		case "topic":
			clusterType = service.ClusterByTopic
		default:
			clusterType = service.ClusterBySender
		}

		clusters := h.clusteringService.ClusterResults(results, clusterType)
		response["clusters"] = clusters
		response["cluster_type"] = clusterTypeStr
	}

	// Add AI summary if enabled
	if enableSummary && h.summaryService != nil && len(results) > 0 {
		summary, err := h.summaryService.GenerateSummary(c.Request.Context(), results, query)
		if err != nil {
			// Fallback to quick summary if AI fails
			h.logger.Warn("AI summary generation failed, using quick summary",
				logger.Error(err),
			)
			summary = h.summaryService.GenerateQuickSummary(results)
		}
		response["summary"] = summary
	}

	c.JSON(http.StatusOK, response)
}
