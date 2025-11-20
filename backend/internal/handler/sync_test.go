package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"echomind.com/backend/internal/handler"
	"echomind.com/backend/pkg/imap"
	clientimap "github.com/emersion/go-imap/client"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// MockFetcher (from service_test) to satisfy handler.EmailFetcher
type MockFetcher struct {
	Results []imap.EmailData
	Err     error
}

func (m *MockFetcher) FetchEmails(c *clientimap.Client, mailbox string, limit int) ([]imap.EmailData, error) {
	return m.Results, m.Err
}

func TestSyncHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// Mock DB and Client (not directly used by handler, but passed to service)
	mockDB := &gorm.DB{}
	mockIMAPClient := &clientimap.Client{}

	// Create a mock fetcher to satisfy the handler's dependencies
	mockFetcher := &MockFetcher{}

	// Create an instance of the handler with mocks (passing nil for asynq client)
	syncHandler := handler.NewSyncHandler(mockDB, mockIMAPClient, mockFetcher, nil)

	// Register the handler route
	r.POST("/api/v1/sync", syncHandler.SyncEmails)

	// Create a test request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/sync", nil)
	r.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, w.Code)
	
	expectedBody := gin.H{"message": "Sync initiated successfully"}
	jsonBody, _ := json.Marshal(expectedBody)
	assert.Equal(t, string(jsonBody), w.Body.String())
}
