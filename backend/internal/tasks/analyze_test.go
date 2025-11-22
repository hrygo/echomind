package tasks

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/hrygo/echomind/internal/model"
	"github.com/hrygo/echomind/pkg/ai"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// MockSummarizer implements the Summarizer interface for testing.
type MockSummarizer struct {
	SummaryResult   ai.AnalysisResult
	SentimentResult ai.SentimentResult
	SummaryError    error
	SentimentError  error
	CallCount       int
}

func (m *MockSummarizer) GenerateSummary(ctx context.Context, text string) (ai.AnalysisResult, error) {
	m.CallCount++
	return m.SummaryResult, m.SummaryError
}

func (m *MockSummarizer) AnalyzeSentiment(ctx context.Context, text string) (ai.SentimentResult, error) {
	return m.SentimentResult, m.SentimentError
}

// MockEmbeddingProvider implements EmbeddingProvider for testing.
type MockEmbeddingProvider struct {
	EmbedResult      []float32
	EmbedBatchResult [][]float32
	EmbedError       error
}

func (m *MockEmbeddingProvider) Embed(ctx context.Context, text string) ([]float32, error) {
	return m.EmbedResult, m.EmbedError
}

func (m *MockEmbeddingProvider) EmbedBatch(ctx context.Context, texts []string) ([][]float32, error) {
	// Return dummy embeddings for each text
	if m.EmbedBatchResult != nil {
		return m.EmbedBatchResult, m.EmbedError
	}
	// Default behavior: return empty vectors of correct length
	var result [][]float32
	for range texts {
		result = append(result, []float32{0.1, 0.2, 0.3})
	}
	return result, m.EmbedError
}

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	if err := db.AutoMigrate(&model.Email{}, &model.Contact{}, &model.EmailEmbedding{}); err != nil {
		t.Fatalf("Failed to auto migrate database: %v", err)
	}
	return db
}

func TestUpdateContactStats_NewContact(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	userID := uuid.New()
	emailAddress := "new@example.com"
	emailSentiment := "Positive"
	interactedAt := time.Now()

	err := updateContactStats(ctx, db, userID, emailAddress, emailSentiment, interactedAt)
	assert.NoError(t, err)

	var contact model.Contact
	err = db.Where("user_id = ? AND email = ?", userID, emailAddress).First(&contact).Error
	assert.NoError(t, err)
	assert.Equal(t, emailAddress, contact.Email)
	assert.Equal(t, 1, contact.InteractionCount)
	assert.Equal(t, 1.0, contact.AvgSentiment)
	assert.WithinDuration(t, interactedAt, contact.LastInteractedAt, time.Second)
}

func TestUpdateContactStats_ExistingContact(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	userID := uuid.New()
	emailAddress := "existing@example.com"
	interactedAt := time.Now().Add(-24 * time.Hour)

	// Create initial contact
	initialContact := model.Contact{
		ID:               uuid.New(),
		UserID:           userID,
		Email:            emailAddress,
		Name:             emailAddress,
		InteractionCount: 1,
		LastInteractedAt: interactedAt,
		AvgSentiment:     0.5,
	}
	db.Create(&initialContact)

	// Update with new interaction
	newInteractedAt := time.Now()
	newSentiment := "Negative"
	err := updateContactStats(ctx, db, userID, emailAddress, newSentiment, newInteractedAt)
	assert.NoError(t, err)

	var updatedContact model.Contact
	err = db.Where("user_id = ? AND email = ?", userID, emailAddress).First(&updatedContact).Error
	assert.NoError(t, err)
	assert.Equal(t, 2, updatedContact.InteractionCount)
	assert.Equal(t, (0.5*1+-1.0)/2, updatedContact.AvgSentiment)
	assert.WithinDuration(t, newInteractedAt, updatedContact.LastInteractedAt, time.Second)
}

func TestHandleEmailAnalyzeTask(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	userID := uuid.New()
	emailID := uuid.New()
	senderEmail := "test-sender@example.com"
	emailDate := time.Now()

	// Create a mock email
	email := model.Email{
		ID:        emailID,
		UserID:    userID,
		MessageID: "<test-message-id>",
		Subject:   "Test Subject",
		Sender:    senderEmail,
		Date:      emailDate,
		BodyText:  "This is a test email body with positive sentiment.",
	}
	db.Create(&email)

	// Mock the summarizer to return specific analysis results
	mockSummarizer := &MockSummarizer{
		SummaryResult: ai.AnalysisResult{
			Summary:     "A positive test summary.",
			Category:    "Personal",
			Sentiment:   "Positive",
			Urgency:     "Low",
			ActionItems: []string{"Reply to sender"},
		},
	}

	mockEmbedder := &MockEmbeddingProvider{}

	// Create the task payload
	payload, _ := json.Marshal(EmailAnalyzePayload{EmailID: emailID, UserID: userID})
	task := asynq.NewTask(TypeEmailAnalyze, payload)

	// Handle the task
	err := HandleEmailAnalyzeTask(ctx, task, db, mockSummarizer, mockEmbedder, 1000)
	assert.NoError(t, err)

	// Verify email was updated
	var updatedEmail model.Email
	db.First(&updatedEmail, "id = ?", emailID)
	assert.Equal(t, "A positive test summary.", updatedEmail.Summary)
	assert.Equal(t, "Personal", updatedEmail.Category)
	assert.Equal(t, "Positive", updatedEmail.Sentiment)
	assert.Equal(t, "Low", updatedEmail.Urgency)

	// Verify contact was updated
	var contact model.Contact
	err = db.Where("user_id = ? AND email = ?", userID, senderEmail).First(&contact).Error
	assert.NoError(t, err)
	assert.Equal(t, 1, contact.InteractionCount)
	assert.Equal(t, 1.0, contact.AvgSentiment)
	assert.WithinDuration(t, emailDate, contact.LastInteractedAt, time.Second)

	// Verify embeddings were saved
	var embeddings []model.EmailEmbedding
	err = db.Where("email_id = ?", emailID).Find(&embeddings).Error
	assert.NoError(t, err)
	assert.NotEmpty(t, embeddings)
}

func TestHandleEmailAnalyzeTask_Spam(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	userID := uuid.New()
	emailID := uuid.New()
	senderEmail := "spammer@example.com"
	emailDate := time.Now()

	// Create a mock spam email
	email := model.Email{
		ID:        emailID,
		UserID:    userID,
		MessageID: "<spam-message-id>",
		Subject:   "Unsubscribe from our newsletter",
		Sender:    senderEmail,
		Date:      emailDate,
		BodyText:  "This is a spam email.",
	}
	db.Create(&email)

	// Mock the summarizer
	mockSummarizer := &MockSummarizer{}
	mockEmbedder := &MockEmbeddingProvider{}

	// Create the task payload
	payload, _ := json.Marshal(EmailAnalyzePayload{EmailID: emailID, UserID: userID})
	task := asynq.NewTask(TypeEmailAnalyze, payload)

	// Handle the task
	err := HandleEmailAnalyzeTask(ctx, task, db, mockSummarizer, mockEmbedder, 1000)
	assert.NoError(t, err)

	// Verify email was updated as spam
	var updatedEmail model.Email
	db.First(&updatedEmail, "id = ?", emailID)
	assert.Equal(t, "Spam", updatedEmail.Category)
	assert.Equal(t, "Neutral", updatedEmail.Sentiment)
	assert.Contains(t, updatedEmail.Summary, "Auto-detected as spam")
	assert.Equal(t, "Low", updatedEmail.Urgency)

	// Verify Summarizer was NOT called
	assert.Equal(t, 0, mockSummarizer.CallCount)
}
