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
	"github.com/hrygo/echomind/internal/middleware"
	"github.com/hrygo/echomind/internal/model"
	"github.com/hrygo/echomind/internal/service"
	"github.com/hrygo/echomind/pkg/imap"
	"github.com/hrygo/echomind/pkg/utils"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
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

func setupSyncTest(t *testing.T) (*handler.SyncHandler, *gorm.DB, *configs.Config) {
	gin.SetMode(gin.TestMode)

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to db: %v", err)
	}
	if err := db.AutoMigrate(&model.Email{}, &model.Contact{}, &model.EmailAccount{}); err != nil {
		t.Fatalf("Failed to auto migrate database: %v", err)
	}

	mockFetcher := &MockFetcher{}
	mockContactService := service.NewContactService(db)
	mockAsynqClient := &asynq.Client{}
	mockConfig := &configs.Config{Security: configs.SecurityConfig{EncryptionKey: "d2f4e23a4b5016b994844b91c48a92c1439bbf17b91a37e4a49ab39c3dbee75f"}}
	mockAccountService := service.NewAccountService(db, &mockConfig.Security)

	mockIMAPClient := &MockIMAPClient{
		DialAndLoginFunc: func(addr, username, password string) (*clientimap.Client, error) {
			return &clientimap.Client{}, nil 
		},
		CloseFunc: func(c *clientimap.Client) { },
	}
	
	syncService := service.NewSyncService(db, mockIMAPClient, mockFetcher, mockAsynqClient, mockContactService, mockAccountService, mockConfig, zap.NewNop().Sugar())
	return handler.NewSyncHandler(syncService), db, mockConfig
}

func TestSyncHandler_Success(t *testing.T) {
	syncHandler, db, mockConfig := setupSyncTest(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodPost, "/api/v1/sync", nil)

	userID := uuid.New()
	c.Set(middleware.ContextUserIDKey, userID)

	encryptionKeyBytes, _ := hex.DecodeString(mockConfig.Security.EncryptionKey)
	encryptedPassword, _ := utils.Encrypt("testpassword", encryptionKeyBytes)

	mockAccount := model.EmailAccount{
		ID:                uuid.New(),
		UserID:            userID,
		Email:             "test@example.com",
		ServerAddress:     "imap.test.com",
		ServerPort:        993,
		Username:          "test@example.com",
		EncryptedPassword: encryptedPassword,
		IsConnected:       true,
	}
	db.Create(&mockAccount)

	syncHandler.SyncEmails(c)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Sync initiated successfully", response["message"])
}

func TestSyncHandler_NoAccount(t *testing.T) {
	syncHandler, _, _ := setupSyncTest(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodPost, "/api/v1/sync", nil)

	userID := uuid.New()
	c.Set(middleware.ContextUserIDKey, userID)

	// Do not create an account for this user

	syncHandler.SyncEmails(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Please configure your email account in Settings first.", response["error"])
}