package spam

import (
	"strings"

	"github.com/hrygo/echomind/internal/model"
)

// Filter defines the interface for spam filtering.
type Filter interface {
	// IsSpam checks if an email is spam.
	// Returns true if spam, along with a reason.
	IsSpam(email *model.Email) (bool, string)
}

// RuleBasedFilter implements a simple rule-based spam filter.
type RuleBasedFilter struct {
	keywords []string
}

// NewRuleBasedFilter creates a new RuleBasedFilter with default rules.
func NewRuleBasedFilter() *RuleBasedFilter {
	return &RuleBasedFilter{
		keywords: []string{
			"unsubscribe",
			"promotion",
			"marketing",
			"verify your email",
			"no-reply",
			"click here",
			"limited time offer",
		},
	}
}

// IsSpam checks if an email is spam based on keywords in subject or body.
func (f *RuleBasedFilter) IsSpam(email *model.Email) (bool, string) {
	// Check Subject
	subjectLower := strings.ToLower(email.Subject)
	for _, keyword := range f.keywords {
		if strings.Contains(subjectLower, keyword) {
			return true, "Subject contains spam keyword: " + keyword
		}
	}

	// Check Body (Snippet is usually a good proxy if BodyText is empty or too long)
	// We check BodyText first, then Snippet.
	textToCheck := email.BodyText
	if textToCheck == "" {
		textToCheck = email.Snippet
	}
	textLower := strings.ToLower(textToCheck)

	for _, keyword := range f.keywords {
		if strings.Contains(textLower, keyword) {
			return true, "Body contains spam keyword: " + keyword
		}
	}

	return false, ""
}
