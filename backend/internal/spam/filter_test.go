package spam

import (
	"testing"

	"github.com/hrygo/echomind/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestRuleBasedFilter_IsSpam(t *testing.T) {
	filter := NewRuleBasedFilter()

	tests := []struct {
		name               string
		email              model.Email
		wantSpam           bool
		wantReasonContains string
	}{
		{
			name: "Clean Email",
			email: model.Email{
				Subject:  "Meeting update",
				BodyText: "Hey, let's meet tomorrow.",
			},
			wantSpam: false,
		},
		{
			name: "Spam Subject - Unsubscribe",
			email: model.Email{
				Subject:  "Unsubscribe from our newsletter",
				BodyText: "Some content",
			},
			wantSpam:           true,
			wantReasonContains: "unsubscribe",
		},
		{
			name: "Spam Body - Promotion",
			email: model.Email{
				Subject:  "Hello",
				BodyText: "Check out this amazing promotion now!",
			},
			wantSpam:           true,
			wantReasonContains: "promotion",
		},
		{
			name: "Spam Snippet - Verify",
			email: model.Email{
				Subject: "Action Required",
				Snippet: "Please verify your email address.",
			},
			wantSpam:           true,
			wantReasonContains: "verify your email",
		},
		{
			name: "Case Insensitive",
			email: model.Email{
				Subject: "MARKETING Opportunity",
			},
			wantSpam:           true,
			wantReasonContains: "marketing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSpam, gotReason := filter.IsSpam(&tt.email)
			assert.Equal(t, tt.wantSpam, gotSpam)
			if tt.wantSpam {
				assert.Contains(t, gotReason, tt.wantReasonContains)
			}
		})
	}
}
