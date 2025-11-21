package handler_test

import (
	"encoding/json"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"testing"

	clientimap "github.com/emersion/go-imap/client"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/hrygo/echomind/configs"
	"github.com/hrygo/echomind/internal/handler"
	"github.com/hrygo/echomind/internal/model"
	"github.com/hrygo/echomind/internal/service"
	"github.com/hrygo/echomind/pkg/imap"
	"github.com/hrygo/echomind/pkg/utils"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
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

// MockIMAPClient implements service.IMAPClient
type MockIMAPClient struct {
	DialAndLoginFunc func(addr, username, password string) (*clientimap.Client, error)
	CloseFunc        func(c *clientimap.Client)
}

func (m *MockIMAPClient) DialAndLogin(addr, username, password string) (*clientimap.Client, error) {
	if m.DialAndLoginFunc != nil {
		return m.DialAndLoginFunc(addr, username, password)
	}
	return &clientimap.Client{}, nil // Return a dummy client
}

func (m *MockIMAPClient) Close(c *clientimap.Client) {
	if m.CloseFunc != nil {
		m.CloseFunc(c)
	}
}

func TestSyncHandler(t *testing.T) {
	// Setup Gin in test mode
	gin.SetMode(gin.TestMode)

	// Mock DB and Client (not directly used by handler, but passed to service)
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to db: %v", err)
	}
	db.AutoMigrate(&model.Email{}, &model.Contact{}, &model.EmailAccount{})
	mockDB := db

	// Create a mock fetcher to satisfy the handler's dependencies
	mockFetcher := &MockFetcher{}

	// Create mock contact service
	mockContactService := service.NewContactService(mockDB)

	// Create mock asynq client, account service and config
	mockAsynqClient := &asynq.Client{}
	mockConfig := &configs.Config{Security: configs.SecurityConfig{EncryptionKey: "d2f4e23a4b5016b994844b91c48a92c1439bbf17b91a37e4a49ab39c3dbee75f"}}
	mockAccountService := service.NewAccountService(mockDB, &mockConfig.Security)

	// Create a mock IMAP client that does nothing but return a dummy client
	mockIMAPClient := &MockIMAPClient{
		DialAndLoginFunc: func(addr, username, password string) (*clientimap.Client, error) {
			return &clientimap.Client{}, nil // Simulate successful connection and login
		},
		CloseFunc: func(c *clientimap.Client) { /* do nothing */ },
	}
	syncService := service.NewSyncService(mockDB, mockIMAPClient, mockFetcher, mockAsynqClient, mockContactService, mockAccountService, mockConfig)

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

	// Create a mock EmailAccount for the user
	encryptionKeyBytes, err := hex.DecodeString(mockConfig.Security.EncryptionKey)
	if err != nil {
		t.Fatalf("Failed to decode encryption key: %v", err)
	}
	encryptedPassword, err := utils.Encrypt("testpassword", encryptionKeyBytes)
	if err != nil {
		t.Fatalf("Failed to encrypt password: %v", err)
	}

	mockAccount := model.EmailAccount{
		ID:                uuid.New(),
		UserID:            userID,
		Email:             "test@example.com",
		ServerAddress:     "imap.test.com",
		ServerPort:        993,
		Username:          "test@example.com",
		EncryptedPassword: encryptedPassword, // Use the actually encrypted password
		IsConnected:       true,
	}
	db.Create(&mockAccount)

	// Call the handler directly
	syncHandler.SyncEmails(c)

	// Assert the response
	assert.Equal(t, http.StatusOK, w.Code)

	expectedBody := gin.H{"message": "Sync initiated successfully"}
	jsonBody, _ := json.Marshal(expectedBody)
	assert.Equal(t, string(jsonBody), w.Body.String())
}
