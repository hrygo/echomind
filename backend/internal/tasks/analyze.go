package tasks

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/hrygo/echomind/internal/model"
	"github.com/hrygo/echomind/pkg/ai"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

const (
	TypeEmailAnalyze = "email:analyze"
)

type EmailAnalyzePayload struct {
	EmailID uuid.UUID
	UserID  uuid.UUID
}

// NewEmailAnalyzeTask creates a task to analyze an email for a specific user.
func NewEmailAnalyzeTask(emailID, userID uuid.UUID) (*asynq.Task, error) {
	payload, err := json.Marshal(EmailAnalyzePayload{EmailID: emailID, UserID: userID})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeEmailAnalyze, payload), nil
}

// Summarizer defines the interface for summary generation and sentiment analysis.
type Summarizer interface {
	GenerateSummary(ctx context.Context, text string) (ai.AnalysisResult, error)
	AnalyzeSentiment(ctx context.Context, text string) (ai.SentimentResult, error)
}

// HandleEmailAnalyzeTask handles the email analysis task for a specific user.
func HandleEmailAnalyzeTask(ctx context.Context, t *asynq.Task, db *gorm.DB, summarizer Summarizer) error {
	var p EmailAnalyzePayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	log.Printf("Analyzing email ID: %s for User ID: %s", p.EmailID, p.UserID)

	// 1. Fetch Email from DB, ensure it belongs to the user
	var email model.Email
	if err := db.WithContext(ctx).Where("id = ? AND user_id = ?", p.EmailID, p.UserID).First(&email).Error; err != nil {
		return fmt.Errorf("email %s not found for user %s: %v", p.EmailID, p.UserID, err)
	}

	// 2. Generate Summary and Analysis
	// Use BodyText or fallback to Snippet/Subject
	textToAnalyze := email.BodyText
	if textToAnalyze == "" {
		textToAnalyze = email.Snippet // Fallback
	}

	analysis, err := summarizer.GenerateSummary(ctx, textToAnalyze)
	if err != nil {
		return fmt.Errorf("failed to generate analysis for email %s (user %s): %v", p.EmailID, p.UserID, err)
	}

	// 3. Update Email fields
	email.Summary = analysis.Summary
	email.Category = analysis.Category
	email.Sentiment = analysis.Sentiment
	email.Urgency = analysis.Urgency
	email.ActionItems = datatypes.JSON(jsonRaw(analysis.ActionItems))

	// 4. Update Email, ensure it belongs to the user
	if err := db.WithContext(ctx).Where("user_id = ?", p.UserID).Save(&email).Error; err != nil {
		return fmt.Errorf("failed to save analysis for email %s (user %s): %v", p.EmailID, p.UserID, err)
	}

	log.Printf("Analysis complete for email %s (user %s). Category: %s, Sentiment: %s", p.EmailID, p.UserID, email.Category, email.Sentiment)

	// 5. Update Contact Statistics for the sender
	if err := updateContactStats(ctx, db, p.UserID, email.Sender, email.Sentiment, email.Date); err != nil {
		log.Printf("Warning: Failed to update contact stats for sender %s: %v", email.Sender, err)
		// Do not return error, as email analysis is complete, contact update can be retried or ignored
	}

	return nil
}

func jsonRaw(v interface{}) []byte {
	b, _ := json.Marshal(v)
	return b
}

// updateContactStats updates or creates a contact and aggregates its statistics.
func updateContactStats(ctx context.Context, db *gorm.DB, userID uuid.UUID, emailAddress string, emailSentiment string, interactedAt time.Time) error {
	var contact model.Contact
	err := db.WithContext(ctx).Where("user_id = ? AND email = ?", userID, emailAddress).First(&contact).Error

	sentimentValue := sentimentToFloat(emailSentiment)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Contact not found, create a new one
			contact = model.Contact{
				ID:               uuid.New(),
				UserID:           userID,
				Email:            emailAddress,
				Name:             emailAddress, // Default name to email address
				InteractionCount: 1,
				LastInteractedAt: interactedAt,
				AvgSentiment:     sentimentValue,
			}
			if createErr := db.WithContext(ctx).Create(&contact).Error; createErr != nil {
				return fmt.Errorf("failed to create contact %s for user %s: %v", emailAddress, userID, createErr)
			}
		} else {
			return fmt.Errorf("failed to query contact %s for user %s: %v", emailAddress, userID, err)
		}
	} else {
		// Contact found, update it
		contact.InteractionCount++
		contact.LastInteractedAt = interactedAt

		// Calculate new average sentiment
		// New AvgSentiment = (Old AvgSentiment * (InteractionCount - 1) + Current Sentiment) / InteractionCount
		contact.AvgSentiment = ((contact.AvgSentiment * float64(contact.InteractionCount-1)) + sentimentValue) / float64(contact.InteractionCount)

		if updateErr := db.WithContext(ctx).Save(&contact).Error; updateErr != nil {
			return fmt.Errorf("failed to update contact %s for user %s: %v", emailAddress, userID, updateErr)
		}
	}
	return nil
}

// sentimentToFloat converts sentiment string to a float value.
func sentimentToFloat(sentiment string) float64 {
	switch sentiment {
	case "Positive":
		return 1.0
	case "Neutral":
		return 0.0
	case "Negative":
		return -1.0
	default:
		return 0.0 // Default to neutral if sentiment is unknown
	}
}
