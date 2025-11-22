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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// MockSearcher implements handler.Searcher
type MockSearcher struct {
	mock.Mock
}

func (m *MockSearcher) Search(ctx context.Context, userID uuid.UUID, query string, limit int) ([]service.SearchResult, error) {
	args := m.Called(ctx, userID, query, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]service.SearchResult), args.Error(1)
}

func TestSearchHandler_Search(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop().Sugar()

	t.Run("Success", func(t *testing.T) {
		mockSearcher := new(MockSearcher)
		h := handler.NewSearchHandler(mockSearcher, logger)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		userID := uuid.New()
		c.Set(middleware.ContextUserIDKey, userID)
		c.Request = httptest.NewRequest("GET", "/api/v1/search?q=project&limit=5", nil)

		expectedResults := []service.SearchResult{
			{EmailID: uuid.New(), Subject: "Test", Score: 0.9},
		}

		mockSearcher.On("Search", mock.Anything, userID, "project", 5).Return(expectedResults, nil)

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
		h := handler.NewSearchHandler(mockSearcher, logger)

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
		h := handler.NewSearchHandler(mockSearcher, logger)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		userID := uuid.New()
		c.Set(middleware.ContextUserIDKey, userID)
		c.Request = httptest.NewRequest("GET", "/api/v1/search?q=error", nil)

		mockSearcher.On("Search", mock.Anything, userID, "error", 10).Return(nil, errors.New("db error"))

		h.Search(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
	
	t.Run("Unauthorized", func(t *testing.T) {
		mockSearcher := new(MockSearcher)
		h := handler.NewSearchHandler(mockSearcher, logger)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		// No user_id set

		c.Request = httptest.NewRequest("GET", "/api/v1/search?q=test", nil)

		h.Search(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
