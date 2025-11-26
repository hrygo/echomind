package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hrygo/echomind/internal/handler"
	"github.com/hrygo/echomind/internal/middleware"
	"github.com/hrygo/echomind/internal/service"
	"github.com/hrygo/echomind/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSearcher implements handler.Searcher
type MockSearcher struct {
	mock.Mock
}

func (m *MockSearcher) Search(ctx context.Context, userID uuid.UUID, query string, filters service.SearchFilters, limit int) ([]service.SearchResult, error) {
	args := m.Called(ctx, userID, query, filters, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]service.SearchResult), args.Error(1)
}

// MockSearchClusterer implements handler.SearchClusterer
type MockSearchClusterer struct {
	mock.Mock
}

func (m *MockSearchClusterer) ClusterResults(results []service.SearchResult, clusterType service.ClusterType) []service.SearchCluster {
	args := m.Called(results, clusterType)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).([]service.SearchCluster)
}

// MockSearchSummarizer implements handler.SearchSummarizer
type MockSearchSummarizer struct {
	mock.Mock
}

func (m *MockSearchSummarizer) GenerateSummary(ctx context.Context, results []service.SearchResult, query string) (*service.SearchResultsSummary, error) {
	args := m.Called(ctx, results, query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.SearchResultsSummary), args.Error(1)
}

func (m *MockSearchSummarizer) GenerateQuickSummary(results []service.SearchResult) *service.SearchResultsSummary {
	args := m.Called(results)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*service.SearchResultsSummary)
}

// NoopLogger implements logger.Logger for testing
type NoopLogger struct{}

func (n *NoopLogger) Debug(msg string, fields ...logger.Field)                             {}
func (n *NoopLogger) Info(msg string, fields ...logger.Field)                              {}
func (n *NoopLogger) Warn(msg string, fields ...logger.Field)                              {}
func (n *NoopLogger) Error(msg string, fields ...logger.Field)                             {}
func (n *NoopLogger) Fatal(msg string, fields ...logger.Field)                             {}
func (n *NoopLogger) DebugContext(ctx context.Context, msg string, fields ...logger.Field) {}
func (n *NoopLogger) InfoContext(ctx context.Context, msg string, fields ...logger.Field)  {}
func (n *NoopLogger) WarnContext(ctx context.Context, msg string, fields ...logger.Field)  {}
func (n *NoopLogger) ErrorContext(ctx context.Context, msg string, fields ...logger.Field) {}
func (n *NoopLogger) FatalContext(ctx context.Context, msg string, fields ...logger.Field) {}
func (n *NoopLogger) With(fields ...logger.Field) logger.Logger                            { return n }
func (n *NoopLogger) SetLevel(level logger.Level)                                          {}
func (n *NoopLogger) GetLevel() logger.Level                                               { return logger.DebugLevel }

func TestSearchHandler_Search(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testLogger := &NoopLogger{}

	t.Run("Success", func(t *testing.T) {
		mockSearcher := new(MockSearcher)
		mockClusterer := new(MockSearchClusterer)
		mockSummarizer := new(MockSearchSummarizer)
		h := handler.NewSearchHandler(mockSearcher, mockClusterer, mockSummarizer, testLogger)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		userID := uuid.New()
		c.Set(middleware.ContextUserIDKey, userID)
		c.Request = httptest.NewRequest("GET", "/api/v1/search?q=project&limit=5", nil)

		expectedResults := []service.SearchResult{
			{EmailID: uuid.New(), Subject: "Test", Score: 0.9},
		}

		mockSearcher.On("Search", mock.Anything, userID, "project", service.SearchFilters{}, 5).Return(expectedResults, nil)

		h.Search(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "project", response["query"])
		assert.Equal(t, float64(1), response["count"])
	})

	t.Run("Missing Query", func(t *testing.T) {
		mockSearcher := new(MockSearcher)
		mockClusterer := new(MockSearchClusterer)
		mockSummarizer := new(MockSearchSummarizer)
		h := handler.NewSearchHandler(mockSearcher, mockClusterer, mockSummarizer, testLogger)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		userID := uuid.New()
		c.Set(middleware.ContextUserIDKey, userID)
		c.Request = httptest.NewRequest("GET", "/api/v1/search", nil) // No query

		h.Search(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Service Error", func(t *testing.T) {
		mockSearcher := new(MockSearcher)
		mockClusterer := new(MockSearchClusterer)
		mockSummarizer := new(MockSearchSummarizer)
		h := handler.NewSearchHandler(mockSearcher, mockClusterer, mockSummarizer, testLogger)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		userID := uuid.New()
		c.Set(middleware.ContextUserIDKey, userID)
		c.Request = httptest.NewRequest("GET", "/api/v1/search?q=error", nil)

		mockSearcher.On("Search", mock.Anything, userID, "error", service.SearchFilters{}, 10).Return(nil, errors.New("db error"))

		h.Search(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		mockSearcher := new(MockSearcher)
		mockClusterer := new(MockSearchClusterer)
		mockSummarizer := new(MockSearchSummarizer)
		h := handler.NewSearchHandler(mockSearcher, mockClusterer, mockSummarizer, testLogger)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		// No user_id set

		c.Request = httptest.NewRequest("GET", "/api/v1/search?q=test", nil)

		h.Search(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
