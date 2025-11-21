package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	clientimap "github.com/emersion/go-imap/client"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hrygo/echomind/internal/handler"
	"github.com/hrygo/echomind/internal/service"
	"github.com/hrygo/echomind/pkg/imap"
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
	// Setup Gin in test mode
	gin.SetMode(gin.TestMode)

	// Mock DB and Client (not directly used by handler, but passed to service)
	mockDB := &gorm.DB{}
	mockIMAPClient := &clientimap.Client{}

	// Create a mock fetcher to satisfy the handler's dependencies
	mockFetcher := &MockFetcher{}

	// Create mock contact service
	mockContactService := service.NewContactService(mockDB)

	// Create sync service with all dependencies
	syncService := service.NewSyncService(mockDB, mockIMAPClient, mockFetcher, nil, mockContactService)

	// Create an instance of the handler with the service
	syncHandler := handler.NewSyncHandler(syncService)

	// Create a Gin context and recorder
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Create a test request
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/sync", nil)
	c.Request = req

	// Set a mock UserID in the context to simulate authentication
	userID := uuid.New()
	c.Set("userID", userID)

	// Call the handler directly
	syncHandler.SyncEmails(c)

	// Assert the response
	assert.Equal(t, http.StatusOK, w.Code)

	expectedBody := gin.H{"message": "Sync initiated successfully"}
	jsonBody, _ := json.Marshal(expectedBody)
	assert.Equal(t, string(jsonBody), w.Body.String())
}
