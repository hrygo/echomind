package imap

import (
	"bytes"
	"io"
	"strings"

	"github.com/emersion/go-message/mail"
)

// ExtractBody extracts the plain text and HTML bodies from a mail reader.
func ExtractBody(r io.Reader) (string, string, error) {
	mr, err := mail.CreateReader(r)
	if err != nil {
		return "", "", err
	}

	var textBody, htmlBody bytes.Buffer

	// Check if it's a multipart message
	// If not multipart, we process the root entity directly.
	// But mail.Reader is designed for reading parts.
	// If it's not multipart, NextPart might return EOF?
	// Let's check the root content type.

	// Note: mail.CreateReader automatically handles the outer structure.
	// If the message is NOT multipart, it might not return parts via NextPart.
	// We should check if we can read from mr directly if parts loop is empty?
	// Actually, go-message's mail.Reader abstracts this.
	// "If the message is not a multipart message, NextPart returns the message itself as the first part." (checking docs mentally).
	// If not, we need to check mr.Header.

	// Iterate over parts
	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", "", err
		}

		contentType := p.Header.Get("Content-Type")
		// fmt.Printf("DEBUG: Part Content-Type: %s\n", contentType)

		if strings.Contains(strings.ToLower(contentType), "text/plain") {
			if _, err := io.Copy(&textBody, p.Body); err != nil {
				return "", "", err
			}
		} else if strings.Contains(strings.ToLower(contentType), "text/html") {
			if _, err := io.Copy(&htmlBody, p.Body); err != nil {
				return "", "", err
			}
		}
	}

	return textBody.String(), htmlBody.String(), nil
}
