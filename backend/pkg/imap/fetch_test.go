package imap

import (
	"testing"
	"time"

	"github.com/emersion/go-imap/server"
)

func TestFetchEmails(t *testing.T) {
	// 1. Setup Mock Server with MockBackend
	be := &MockBackend{}
	s := server.New(be)
	s.Addr = "127.0.0.1:3001" 
	s.AllowInsecureAuth = true
	
go func() {
		_ = s.ListenAndServe()
	}()
	defer s.Close()
	time.Sleep(100 * time.Millisecond)

	// 2. Connect
	c, err := Connect("127.0.0.1:3001", "user", "pass", false)
	if err != nil {
		t.Fatalf("Connect failed: %v", err)
	}
	defer c.Logout()

	// 3. Call FetchEmails
	emails, err := FetchEmails(c, "INBOX", 10)
	if err != nil {
		t.Fatalf("FetchEmails failed: %v", err)
	}

	if len(emails) != 1 {
		t.Fatalf("Expected 1 email, got %d", len(emails))
	}
	if emails[0].Subject != "Test Subject" {
		t.Errorf("Expected subject 'Test Subject', got '%s'", emails[0].Subject)
	}
	// Note: ExtractBody logic checks for "text/plain" or "text/html".
	// The mock body "Content-Type: text/plain\r\n\r\nThis is a test body." should be parsed.
	// However, go-message parsing might be strict about headers.
	// Let's check if we got the body.
	if emails[0].BodyText != "This is a test body." {
		t.Errorf("Expected BodyText 'This is a test body.', got '%s'", emails[0].BodyText)
	}
}
