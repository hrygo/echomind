package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hrygo/echomind/internal/model"
	"github.com/hrygo/echomind/internal/service"
	"github.com/hrygo/echomind/pkg/imap"
	"github.com/emersion/go-imap/client"
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

func TestSyncEmails(t *testing.T) {
	// 1. Setup DB
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to db: %v", err)
	}
	db.AutoMigrate(&model.Email{}, &model.Contact{}) // Also migrate Contact model

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
	syncService := service.NewSyncService(db, nil, fetcher, nil, mockContactService)

	// 5. Run Sync
	userID := uuid.New() // Generate a new UserID for the test
	ctx := context.Background()
	err = syncService.SyncEmails(ctx, userID)
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
