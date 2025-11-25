package service_test

import (
	"context"
	"encoding/hex"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hrygo/echomind/configs"
	"github.com/hrygo/echomind/internal/event"
	"github.com/hrygo/echomind/internal/listener"
	"github.com/hrygo/echomind/internal/model"
	"github.com/hrygo/echomind/internal/repository"
	"github.com/hrygo/echomind/internal/service"
	"github.com/hrygo/echomind/pkg/event/bus"
	"github.com/hrygo/echomind/pkg/imap"
	"github.com/hrygo/echomind/pkg/logger"
	"github.com/hrygo/echomind/pkg/utils"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// }

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

func TestSyncEmails(t *testing.T) {
	// 1. Setup DB
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to db: %v", err)
	}
	if err := db.AutoMigrate(&model.Email{}, &model.Contact{}, &model.EmailAccount{}); err != nil {
		t.Fatalf("Failed to auto migrate database: %v", err)
	}
	// 2. Setup Mock Data
	now := time.Now()
	mockData := []imap.EmailData{
		{
			Subject:   "Sync Test",
			Sender:    "Sync Test <sync@test.com>",
			Date:      now,
			MessageID: "<sync@test.com>",
			BodyText:  "Test Body Content",
		},
	}

	// 3. Setup Mock ContactService
	// We need this for the ContactListener
	mockContactService := service.NewContactService(db)

	// 4. Setup SyncService (SUT)
	mockConfig := &configs.Config{Security: configs.SecurityConfig{EncryptionKey: "d2f4e23a4b5016b994844b91c48a92c1439bbf17b91a37e4a49ab39c3dbee75f"}}
	mockAccountService := service.NewAccountService(db, &mockConfig.Security)

	// Create Mock Connector and Session
	mockSession := &MockIMAPSession{
		FetchEmailsFunc: func(mailbox string, limit int) ([]imap.EmailData, error) {
			return mockData, nil
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
	emailRepo := repository.NewEmailRepository(db)
	accountRepo := repository.NewAccountRepository(db)
	eventBus := bus.New()

	// Wire up ContactListener
	contactListener := listener.NewContactListener(mockContactService, logger.GetDefaultLogger())
	eventBus.Subscribe(event.EmailSyncedEventName, contactListener)

	ingestor := service.NewEmailIngestor(emailRepo, logger.GetDefaultLogger())

	syncService := service.NewSyncService(accountRepo, mockConnector, ingestor, eventBus, mockAccountService, mockConfig, logger.GetDefaultLogger())

	// 5. Create a mock EmailAccount for the user
	userID := uuid.New() // Generate a new UserID for the test

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
		UserID:            &userID,
		Email:             "test@example.com",
		ServerAddress:     "imap.test.com",
		ServerPort:        993,
		Username:          "test@example.com",
		EncryptedPassword: encryptedPassword, // Use the actually encrypted password
		IsConnected:       true,
	}
	db.Create(&mockAccount)

	ctx := context.Background()
	err = syncService.SyncEmails(ctx, userID, nil, nil)
	if err != nil {
		t.Fatalf("SyncEmails failed: %v", err)
	}

	// 6. Verify DB
	var count int64
	db.Model(&model.Email{}).Where("user_id = ?", userID).Count(&count)
	if count != 1 {
		t.Errorf("Expected 1 email, got %d", count)
	}

	var email model.Email
	db.Where("user_id = ?", userID).First(&email)
	if email.Subject != "Sync Test" {
		t.Errorf("Expected subject 'Sync Test', got '%s'", email.Subject)
	}
	if email.BodyText != "Test Body Content" {
		t.Errorf("Expected body 'Test Body Content', got '%s'", email.BodyText)
	}

	var contact model.Contact
	db.Where("user_id = ? AND email = ?", userID, "sync@test.com").First(&contact)
	if contact.Name != "Sync Test" {
		t.Errorf("Expected contact name 'Sync Test', got '%s'", contact.Name)
	}
	if contact.InteractionCount != 1 {
		t.Errorf("Expected contact interaction count 1, got %d", contact.InteractionCount)
	}
}
