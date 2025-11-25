package model_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hrygo/echomind/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestEmailModel(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	// 1. Test AutoMigrate (Schema Check)
	// This will fail if model.Email is undefined
	if err := db.AutoMigrate(&model.Email{}); err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	// 2. Test Create
	now := time.Now()
	email := model.Email{
		ID:        uuid.New(),
		UserID:    uuid.New(),
		MessageID: "<123@example.com>",
		Subject:   "Test Subject",
		Sender:    "sender@example.com",
		Date:      now,
		Snippet:   "This is a snippet...",
		Sentiment: "Positive",
		Urgency:   "High",
	}

	if err := db.Create(&email).Error; err != nil {
		t.Fatalf("Failed to create email: %v", err)
	}

	// 3. Test Read
	var readEmail model.Email
	if err := db.First(&readEmail, "message_id = ?", "<123@example.com>").Error; err != nil {
		t.Fatalf("Failed to query email: %v", err)
	}

	if readEmail.Subject != "Test Subject" {
		t.Errorf("Expected subject 'Test Subject', got '%s'", readEmail.Subject)
	}
	if readEmail.Sentiment != "Positive" {
		t.Errorf("Expected sentiment 'Positive', got '%s'", readEmail.Sentiment)
	}
	if readEmail.Urgency != "High" {
		t.Errorf("Expected urgency 'High', got '%s'", readEmail.Urgency)
	}
}
