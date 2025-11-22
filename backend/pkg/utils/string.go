package utils

import (
	"regexp"
	"strings"
)

// Slugify converts a string to a slug (lowercase, dashes)
func Slugify(s string) string {
	// Convert to lowercase
	s = strings.ToLower(s)
	// Replace non-alphanumeric chars with dashes
	reg := regexp.MustCompile("[^a-z0-9]+")
	s = reg.ReplaceAllString(s, "-")
	// Trim dashes
	return strings.Trim(s, "-")
}
