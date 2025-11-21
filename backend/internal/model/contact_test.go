package model_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hrygo/echomind/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestContactModel(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(&model.Contact{}); err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	// Test Create
	now := time.Now()
	contact := model.Contact{
		ID:               uuid.New(),
		UserID:           uuid.New(),
		Email:            "test@example.com",
		Name:             "Test User",
		InteractionCount: 1,
		LastInteractedAt: now,
	}

	if err := db.Create(&contact).Error; err != nil {
		t.Fatalf("Failed to create contact: %v", err)
	}

	// Test Read
	var readContact model.Contact
	if err := db.First(&readContact, "email = ?", "test@example.com").Error; err != nil {
		t.Fatalf("Failed to query contact: %v", err)
	}

	if readContact.Name != "Test User" {
		t.Errorf("Expected name 'Test User', got '%s'", readContact.Name)
	}
}
