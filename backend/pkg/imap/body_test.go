package imap

import (
	"strings"
	"testing"
)

func TestExtractBody(t *testing.T) {
	// 1. Construct a raw multipart email (Text + HTML)
	rawEmail := "From: sender@example.com\r\n" +
		"To: recipient@example.com\r\n" +
		"Subject: Test Email\r\n" +
		"Content-Type: multipart/alternative; boundary=boundary123\r\n" +
		"\r\n" +
		"--boundary123\r\n" +
		"Content-Type: text/plain; charset=utf-8\r\n" +
		"\r\n" +
		"This is the plain text body.\r\n" +
		"--boundary123\r\n" +
		"Content-Type: text/html; charset=utf-8\r\n" +
		"\r\n" +
		"<html><body>This is the <b>HTML</b> body.</body></html>\r\n" +
		"--boundary123--\r\n"

	// 2. Parse it using go-message (simulating what we'd do with the IMAP literal)
	r := strings.NewReader(rawEmail)
	
	// 3. Call ExtractBody (SUT)
	textBody, htmlBody, err := ExtractBody(r)
	if err != nil {
		t.Fatalf("ExtractBody failed: %v", err)
	}

	// 4. Assertions
	expectedText := "This is the plain text body."
	if strings.TrimSpace(textBody) != expectedText {
		t.Errorf("Expected text body '%s', got '%s'", expectedText, textBody)
	}

	expectedHTML := "<html><body>This is the <b>HTML</b> body.</body></html>"
	if strings.TrimSpace(htmlBody) != expectedHTML {
		t.Errorf("Expected html body '%s', got '%s'", expectedHTML, htmlBody)
	}
}

func TestExtractBody_Simple(t *testing.T) {
	rawEmail := "Content-Type: text/plain; charset=utf-8\r\n" +
		"\r\n" +
		"Simple body."

	r := strings.NewReader(rawEmail)
	textBody, _, err := ExtractBody(r)
	if err != nil {
		t.Fatalf("ExtractBody failed: %v", err)
	}

	if strings.TrimSpace(textBody) != "Simple body." {
		t.Errorf("Expected 'Simple body.', got '%s'", textBody)
	}
}

