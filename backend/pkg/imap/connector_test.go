package imap

import (
	"testing"
	"time"

	"github.com/emersion/go-imap/backend/memory"
	"github.com/emersion/go-imap/server"
)

func TestConnect(t *testing.T) {
	// 1. Setup Mock Server
	be := memory.New()
	
	s := server.New(be)
	s.Addr = "127.0.0.1:3000"
	s.AllowInsecureAuth = true
	
	// Start server
	go func() {
		_ = s.ListenAndServe()
	}()
	defer s.Close()
	
	// Allow server to start
	time.Sleep(100 * time.Millisecond)

	// 2. Test Connect function
	// We pass useTLS = false for this test
	c, err := Connect("127.0.0.1:3000", "user", "pass", false)
	
	if err == nil {
		t.Fatal("Expected error due to missing user, but got nil")
		c.Logout()
	}

	// Verify it's an auth error, not a network error
	// "Bad username or password" is the standard error from go-imap memory backend
	if err.Error() != "Bad username or password" {
		t.Errorf("Expected 'Bad username or password', got: %v", err)
	}
}