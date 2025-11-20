package service_test

import (
	"testing"
	"time"

	"echomind.com/backend/internal/model"
	"echomind.com/backend/internal/service"
	"echomind.com/backend/pkg/imap"
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
	db.AutoMigrate(&model.Email{})

	// 2. Setup Mock Fetcher
	now := time.Now()
	mockData := []imap.EmailData{
		{
			Subject:   "Sync Test",
			Sender:    "sync@test.com",
			Date:      now,
			MessageID: "<sync@test.com>",
			BodyText:  "Test Body Content",
		},
	}
	fetcher := &MockFetcher{Results: mockData}

	// 3. Run Sync (SUT)
	// We pass nil for client and asynq client (optional)
	err = service.SyncEmails(db, nil, fetcher, nil)
	if err != nil {
		t.Fatalf("SyncEmails failed: %v", err)
	}

	// 4. Verify DB
	var count int64
	db.Model(&model.Email{}).Count(&count)
	if count != 1 {
		t.Errorf("Expected 1 email, got %d", count)
	}

	var email model.Email
	db.First(&email)
	if email.Subject != "Sync Test" {
		t.Errorf("Expected subject 'Sync Test', got '%s'", email.Subject)
	}
	if email.BodyText != "Test Body Content" {
		t.Errorf("Expected body 'Test Body Content', got '%s'", email.BodyText)
	}
}
