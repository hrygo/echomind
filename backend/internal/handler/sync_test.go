package handler_test

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	clientimap "github.com/emersion/go-imap/client"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hrygo/echomind/configs"
	"github.com/hrygo/echomind/internal/handler"
	"github.com/hrygo/echomind/internal/middleware"
	"github.com/hrygo/echomind/internal/model"
	"github.com/hrygo/echomind/internal/repository"
	"github.com/hrygo/echomind/internal/service"
	"github.com/hrygo/echomind/pkg/event/bus"
	"github.com/hrygo/echomind/pkg/imap"
	"github.com/hrygo/echomind/pkg/logger"
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

// MockIMAPSession implements service.IMAPSession
type MockIMAPSession struct {
	LogoutFunc      func() error
	FetchEmailsFunc func(mailbox string, limit int) ([]imap.EmailData, error)
}

func (m *MockIMAPSession) Logout() error {
	if m.LogoutFunc != nil {
		return m.LogoutFunc()
	}
	return nil
}

func (m *MockIMAPSession) FetchEmails(mailbox string, limit int) ([]imap.EmailData, error) {
	if m.FetchEmailsFunc != nil {
		return m.FetchEmailsFunc(mailbox, limit)
	}
	return nil, nil
}

// MockIMAPConnector implements service.IMAPConnector
type MockIMAPConnector struct {
	ConnectFunc func(ctx context.Context, account *model.EmailAccount) (service.IMAPSession, error)
}

func (m *MockIMAPConnector) Connect(ctx context.Context, account *model.EmailAccount) (service.IMAPSession, error) {
	if m.ConnectFunc != nil {
		return m.ConnectFunc(ctx, account)
	}
	return &MockIMAPSession{}, nil
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

	mockConfig := &configs.Config{Security: configs.SecurityConfig{EncryptionKey: "d2f4e23a4b5016b994844b91c48a92c1439bbf17b91a37e4a49ab39c3dbee75f"}}
	mockAccountService := service.NewAccountService(db, &mockConfig.Security)

	// Create Mock Connector and Session
	mockSession := &MockIMAPSession{
		FetchEmailsFunc: func(mailbox string, limit int) ([]imap.EmailData, error) {
			return []imap.EmailData{}, nil
		},
	}
	mockConnector := &MockIMAPConnector{
		ConnectFunc: func(ctx context.Context, account *model.EmailAccount) (service.IMAPSession, error) {
			return mockSession, nil
		},
	}

	// 初始化一个简单的 logger 用于测试
	if err := logger.Init(logger.DevelopmentConfig()); err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}
	// The second logger.Init call is redundant, removing it.
	// if err := logger.Init(logger.DevelopmentConfig()); err != nil {
	// 	t.Fatalf("Failed to initialize logger: %v", err)
	// }
	emailRepo := repository.NewEmailRepository(db)
	accountRepo := repository.NewAccountRepository(db)
	eventBus := bus.New()

	ingestor := service.NewEmailIngestor(emailRepo, logger.GetDefaultLogger())

	syncService := service.NewSyncService(accountRepo, mockConnector, ingestor, eventBus, mockAccountService, mockConfig, logger.GetDefaultLogger())
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
		UserID:            &userID,
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
	_ = json.Unmarshal(w.Body.Bytes(), &response)
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
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Please configure your email account in Settings first.", response["error"])
}
