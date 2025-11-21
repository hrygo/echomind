package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/hrygo/echomind/internal/model"
	"github.com/hrygo/echomind/pkg/ai"
	"github.com/hibiken/asynq"
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
	GenerateSummary(ctx context.Context, text string) (string, error)
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

	// 2. Generate Summary
	// Use BodyText or fallback to Snippet/Subject
	textToAnalyze := email.BodyText
	if textToAnalyze == "" {
		textToAnalyze = email.Snippet // Fallback
	}
	
	summary, err := summarizer.GenerateSummary(ctx, textToAnalyze)
	if err != nil {
		return fmt.Errorf("failed to generate summary for email %s (user %s): %v", p.EmailID, p.UserID, err)
	}

	// 3. Analyze Sentiment
	sentimentResult, err := summarizer.AnalyzeSentiment(ctx, textToAnalyze)
	if err != nil {
		log.Printf("Warning: sentiment analysis failed for email %s (user %s): %v", email.ID, p.UserID, err)
		// Continue without sentiment? Or fail? Let's continue for now.
	} else {
		email.Sentiment = sentimentResult.Sentiment
		email.Urgency = sentimentResult.Urgency
	}

	// 4. Update Email, ensure it belongs to the user
	email.Summary = summary
	if err := db.WithContext(ctx).Where("user_id = ?", p.UserID).Save(&email).Error; err != nil {
		return fmt.Errorf("failed to save analysis for email %s (user %s): %v", p.EmailID, p.UserID, err)
	}
	
	log.Printf("Analysis complete for email %s (user %s). Summary: %s, Sentiment: %s", p.EmailID, p.UserID, summary, email.Sentiment)
	return nil
}