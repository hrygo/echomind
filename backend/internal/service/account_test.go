package service

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/hrygo/echomind/configs"
	"github.com/hrygo/echomind/internal/model"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupAccountTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	_ = db.AutoMigrate(&model.EmailAccount{})
	return db
}

func TestConnectAndSaveAccount_WithSMTP(t *testing.T) {
	db := setupAccountTestDB()

	// Mock Config
	// 32 bytes hex string for AES-256
	mockKey := "000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f"
	config := &configs.SecurityConfig{
		EncryptionKey: mockKey,
	}

	svc := NewAccountService(db, config)
	ctx := context.Background()
	userID := uuid.New()

	input := &model.EmailAccountInput{
		Email:         "mock@test.com",
		Username:      "mock@test.com", // Triggers mock connection success
		Password:      "password123",
		ServerAddress: "imap.test.com",
		ServerPort:    993,
		SMTPServer:    "smtp.test.com",
		SMTPPort:      587,
	}

	// Test Create
	account, err := svc.ConnectAndSaveAccount(ctx, userID, input)
	assert.NoError(t, err)
	assert.NotNil(t, account)
	assert.Equal(t, "smtp.test.com", account.SMTPServer)
	assert.Equal(t, 587, account.SMTPPort)
	assert.Equal(t, "imap.test.com", account.ServerAddress)
	assert.Equal(t, 993, account.ServerPort)

	// Verify in DB
	var savedAccount model.EmailAccount
	err = db.First(&savedAccount, "id = ?", account.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, "smtp.test.com", savedAccount.SMTPServer)
	assert.Equal(t, 587, savedAccount.SMTPPort)

	// Test Update
	input.SMTPServer = "smtp.updated.com"
	input.SMTPPort = 465

	updatedAccount, err := svc.ConnectAndSaveAccount(ctx, userID, input)
	assert.NoError(t, err)
	if err != nil {
		return
	}
	assert.Equal(t, "smtp.updated.com", updatedAccount.SMTPServer)
	assert.Equal(t, 465, updatedAccount.SMTPPort)
	assert.Equal(t, account.ID, updatedAccount.ID) // Should be same ID
}

func TestDisconnectAccount(t *testing.T) {
	db := setupAccountTestDB()
	svc := NewAccountService(db, &configs.SecurityConfig{EncryptionKey: "000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f"})
	ctx := context.Background()
	userID := uuid.New()

	// 1. Create account
	input := &model.EmailAccountInput{
		Email:         "mock@test.com",
		Username:      "mock@test.com",
		Password:      "password123",
		ServerAddress: "imap.test.com",
		ServerPort:    993,
	}
	_, err := svc.ConnectAndSaveAccount(ctx, userID, input)
	assert.NoError(t, err)

	// 2. Disconnect account
	err = svc.DisconnectAccount(ctx, userID)
	assert.NoError(t, err)

	// 3. Verify deletion
	var count int64
	db.Model(&model.EmailAccount{}).Where("user_id = ?", userID).Count(&count)
	assert.Equal(t, int64(0), count)

	// 4. Disconnect again (should fail with RecordNotFound)
	err = svc.DisconnectAccount(ctx, userID)
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}
