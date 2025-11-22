package service_test

import (
	"context"
	"encoding/hex"
	"testing"
	"time"

	"github.com/emersion/go-imap/client"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/hrygo/echomind/configs"
	"github.com/hrygo/echomind/internal/model"
	"github.com/hrygo/echomind/internal/service"
	"github.com/hrygo/echomind/pkg/imap"
	"github.com/hrygo/echomind/pkg/utils"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// MockFetcher implements service.EmailFetcher
type MockFetcher struct {
	Results []imap.EmailData
	Err     error
}

func (m *MockFetcher) FetchEmails(c *client.Client, mailbox string, limit int) ([]imap.EmailData, error) {
	return m.Results, m.Err
}

// MockIMAPClient implements service.IMAPClient
type MockIMAPClient struct {
	DialAndLoginFunc func(addr, username, password string) (*client.Client, error)
	CloseFunc        func(c *client.Client)
}

func (m *MockIMAPClient) DialAndLogin(addr, username, password string) (*client.Client, error) {
	if m.DialAndLoginFunc != nil {
		return m.DialAndLoginFunc(addr, username, password)
	}
	return &client.Client{}, nil // Return a dummy client
}

func (m *MockIMAPClient) Close(c *client.Client) {
	if m.CloseFunc != nil {
		m.CloseFunc(c)
	}
}

// MockAsynqClient implements asynqClientInterface for testing.
type MockAsynqClient struct{
	EnqueueFunc func(task *asynq.Task, opts ...asynq.Option) (*asynq.TaskInfo, error)
}

func (m *MockAsynqClient) Enqueue(task *asynq.Task, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	if m.EnqueueFunc != nil {
		return m.EnqueueFunc(task, opts...)
	}
	return &asynq.TaskInfo{}, nil
}

func (m *MockAsynqClient) Close() error {
	return nil
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
	// 2. Setup Mock Fetcher
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
	fetcher := &MockFetcher{Results: mockData}

	// 3. Setup Mock ContactService
	// For this test, we don't need to mock its behavior extensively, just provide an instance.
	mockContactService := service.NewContactService(db) // Assuming NewContactService only needs db

	// 4. Setup SyncService (SUT)
	// We pass nil for client and asynq client for this unit test
	mockAsynqClient := &MockAsynqClient{} // Use MockAsynqClient
	mockConfig := &configs.Config{Security: configs.SecurityConfig{EncryptionKey: "d2f4e23a4b5016b994844b91c48a92c1439bbf17b91a37e4a49ab39c3dbee75f"}}
	mockAccountService := service.NewAccountService(db, &mockConfig.Security)

	// Create a mock IMAP client that does nothing but return a dummy client
	mockIMAPClient := &MockIMAPClient{
		DialAndLoginFunc: func(addr, username, password string) (*client.Client, error) {
			return &client.Client{}, nil // Simulate successful connection and login
		},
		CloseFunc: func(c *client.Client) { /* do nothing */ },
	}
	syncService := service.NewSyncService(db, mockIMAPClient, fetcher, mockAsynqClient, mockContactService, mockAccountService, mockConfig, zap.NewNop().Sugar())

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
