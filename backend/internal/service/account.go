package service

import (
	"context"
	"fmt"
	"time"

	clientimap "github.com/emersion/go-imap/client"
	"github.com/google/uuid"
	"github.com/hrygo/echomind/configs"
	"github.com/hrygo/echomind/internal/model"
	"github.com/hrygo/echomind/pkg/utils"
	"gorm.io/gorm"
)

// AccountService handles operations related to user email accounts.
type AccountService struct {
	db           *gorm.DB
	config       *configs.SecurityConfig
}

// NewAccountService creates a new AccountService.
func NewAccountService(db *gorm.DB, config *configs.SecurityConfig) *AccountService {
	return &AccountService{
		db:           db,
		config:       config,
	}
}

// ConnectAndSaveAccount attempts to connect to an IMAP server with provided credentials,
// encrypts the password if successful, and saves/updates the EmailAccount in the database.
func (s *AccountService) ConnectAndSaveAccount(ctx context.Context, userID uuid.UUID, input *model.EmailAccountInput) (*model.EmailAccount, error) {
	// 1. Test IMAP connection
	if err := s.testIMAPConnection(input.ServerAddress, input.ServerPort, input.Username, input.Password); err != nil {
		return nil, fmt.Errorf("IMAP connection test failed: %w", err)
	}

	// 2. Encrypt password
	encryptedPassword, err := utils.Encrypt(input.Password, []byte(s.config.EncryptionKey))
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt password: %w", err)
	}

	// 3. Prepare EmailAccount model
	account := model.EmailAccount{
		UserID:            userID,
		Email:             input.Email,
		ServerAddress:     input.ServerAddress,
		ServerPort:        input.ServerPort,
		Username:          input.Username,
		EncryptedPassword: encryptedPassword,
		IsConnected:       true,
		LastSyncAt:        nil, // Will be set on first successful sync
		ErrorMessage:      "",
	}

	// 4. Upsert (Create or Update) the account
	var existingAccount model.EmailAccount
	res := s.db.WithContext(ctx).Where("user_id = ?", userID).First(&existingAccount)

	if res.Error == gorm.ErrRecordNotFound {
		// Create new account
		if err := s.db.WithContext(ctx).Create(&account).Error; err != nil {
			return nil, fmt.Errorf("failed to create email account: %w", err)
		}
	} else if res.Error != nil {
		return nil, fmt.Errorf("database error while checking existing account: %w", res.Error)
	} else {
		// Update existing account
		account.ID = existingAccount.ID // Retain existing ID
		if err := s.db.WithContext(ctx).Save(&account).Error; err != nil {
			return nil, fmt.Errorf("failed to update email account: %w", err)
		}
	}

	return &account, nil
}

// GetAccountByUserID retrieves an EmailAccount by UserID.
// It decrypts the password for use, but does NOT return it in the model.
func (s *AccountService) GetAccountByUserID(ctx context.Context, userID uuid.UUID) (*model.EmailAccount, error) {
	var account model.EmailAccount
	if err := s.db.WithContext(ctx).Where("user_id = ?", userID).First(&account).Error; err != nil {
		return nil, fmt.Errorf("email account not found for user %s: %w", userID, err)
	}
	return &account, nil
}

// UpdateAccountStatus updates the connection status and error message for an email account.
func (s *AccountService) UpdateAccountStatus(ctx context.Context, accountID uuid.UUID, isConnected bool, errorMessage string, lastSyncAt *time.Time) error {
	return s.db.WithContext(ctx).Model(&model.EmailAccount{}).Where("id = ?", accountID).Updates(map[string]interface{}{
		"is_connected":  isConnected,
		"error_message": errorMessage,
		"last_sync_at":  lastSyncAt,
	}).Error
}

// testIMAPConnection attempts to establish a basic IMAP connection and login.
func (s *AccountService) testIMAPConnection(server string, port int, username, password string) error {
	addr := fmt.Sprintf("%s:%d", server, port)
	client, err := clientimap.DialTLS(addr, nil)
	if err != nil {
		return fmt.Errorf("dial IMAP server %s failed: %w", addr, err)
	}
	defer client.Close()

	if err := client.Login(username, password); err != nil {
		return fmt.Errorf("IMAP login failed: %w", err)
	}
	return nil
}
