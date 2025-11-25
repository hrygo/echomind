package service

import (
	"context"
	"encoding/hex"
	"fmt"

	clientimap "github.com/emersion/go-imap/client"
	"github.com/hrygo/echomind/configs"
	"github.com/hrygo/echomind/internal/model"
	"github.com/hrygo/echomind/pkg/imap"
	"github.com/hrygo/echomind/pkg/utils"
)

// IMAPSession defines the interface for an authenticated IMAP session.
type IMAPSession interface {
	Logout() error
	FetchEmails(mailbox string, limit int) ([]imap.EmailData, error)
}

// DefaultIMAPSession wraps a go-imap client.
type DefaultIMAPSession struct {
	client *clientimap.Client
}

func (s *DefaultIMAPSession) Logout() error {
	return s.client.Logout()
}

func (s *DefaultIMAPSession) FetchEmails(mailbox string, limit int) ([]imap.EmailData, error) {
	return imap.FetchEmails(s.client, mailbox, limit)
}

// IMAPConnector handles establishing connections to IMAP servers.
type IMAPConnector interface {
	Connect(ctx context.Context, account *model.EmailAccount) (IMAPSession, error)
}

// DefaultIMAPConnector implements IMAPConnector.
type DefaultIMAPConnector struct {
	clientFactory IMAPClient
	config        *configs.Config
}

func NewIMAPConnector(clientFactory IMAPClient, config *configs.Config) *DefaultIMAPConnector {
	return &DefaultIMAPConnector{
		clientFactory: clientFactory,
		config:        config,
	}
}

// Connect establishes an authenticated connection to the IMAP server for the given account.
func (c *DefaultIMAPConnector) Connect(ctx context.Context, account *model.EmailAccount) (IMAPSession, error) {
	// 1. Decrypt password
	encryptionKeyBytes, err := hex.DecodeString(c.config.Security.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode encryption key: %w", err)
	}

	password, err := utils.Decrypt(account.EncryptedPassword, encryptionKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt password: %w", err)
	}

	// 2. Connect to server
	addr := fmt.Sprintf("%s:%d", account.ServerAddress, account.ServerPort)
	client, err := c.clientFactory.DialAndLogin(addr, account.Username, password)
	if err != nil {
		return nil, fmt.Errorf("failed to connect/login to IMAP server: %w", err)
	}

	return &DefaultIMAPSession{client: client}, nil
}
