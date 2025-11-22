package utils

import (
	"regexp"
	"strings"
)

// StripHTML removes HTML tags from a string and returns the plain text.
// This is a regex-based implementation for simplicity and performance.
// For more complex cases, a proper HTML parser (golang.org/x/net/html) should be used.
func StripHTML(htmlContent string) string {
	if htmlContent == "" {
		return ""
	}

	// Regex to match HTML tags
	re := regexp.MustCompile(`<[^>]*>`)
	text := re.ReplaceAllString(htmlContent, " ")

	// Decode common HTML entities (basic set)
	text = strings.ReplaceAll(text, "&nbsp;", " ")
	text = strings.ReplaceAll(text, "&amp;", "&")
	text = strings.ReplaceAll(text, "&lt;", "<")
	text = strings.ReplaceAll(text, "&gt;", ">")
	text = strings.ReplaceAll(text, "&quot;", "\"")
	text = strings.ReplaceAll(text, "&#39;", "'")

	// Collapse multiple spaces into one
	spaceRe := regexp.MustCompile(`\s+`)
	text = spaceRe.ReplaceAllString(text, " ")

	return strings.TrimSpace(text)
}
