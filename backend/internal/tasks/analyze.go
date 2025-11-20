package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"echomind.com/backend/internal/model"
	"echomind.com/backend/pkg/ai"
	"github.com/hibiken/asynq"
	"gorm.io/gorm"
)

const (
	TypeEmailAnalyze = "email:analyze"
)

type EmailAnalyzePayload struct {
	EmailID uint
}

// NewEmailAnalyzeTask creates a task to analyze an email.
func NewEmailAnalyzeTask(emailID uint) (*asynq.Task, error) {
	payload, err := json.Marshal(EmailAnalyzePayload{EmailID: emailID})
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

// HandleEmailAnalyzeTask handles the email analysis task.
func HandleEmailAnalyzeTask(ctx context.Context, t *asynq.Task, db *gorm.DB, summarizer Summarizer) error {
	var p EmailAnalyzePayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	log.Printf("Analyzing email ID: %d", p.EmailID)

	// 1. Fetch Email from DB
	var email model.Email
	if err := db.First(&email, p.EmailID).Error; err != nil {
		return fmt.Errorf("email not found: %v", err)
	}

	// 2. Generate Summary
	// Use BodyText or fallback to Snippet/Subject
	textToAnalyze := email.BodyText
	if textToAnalyze == "" {
		textToAnalyze = email.Snippet // Fallback
	}
	
	summary, err := summarizer.GenerateSummary(ctx, textToAnalyze)
	if err != nil {
		return fmt.Errorf("failed to generate summary: %v", err)
	}

	// 3. Analyze Sentiment
	sentimentResult, err := summarizer.AnalyzeSentiment(ctx, textToAnalyze)
	if err != nil {
		log.Printf("Warning: sentiment analysis failed for email %d: %v", email.ID, err)
		// Continue without sentiment? Or fail? Let's continue for now.
	} else {
		email.Sentiment = sentimentResult.Sentiment
		email.Urgency = sentimentResult.Urgency
	}

	// 4. Update Email
	email.Summary = summary
	if err := db.Save(&email).Error; err != nil {
		return fmt.Errorf("failed to save analysis: %v", err)
	}
	
	log.Printf("Analysis complete for email %d. Summary: %s, Sentiment: %s", p.EmailID, summary, email.Sentiment)
	return nil
}

// Ensure pkg/ai/provider.go's AIProvider is not used directly if Summarizer signature matches.
// Or we can define Summarizer interface here.