package ai

import "context"

// AIProvider defines the interface for AI capabilities.
type AIProvider interface {
	// Summarize generates a structured analysis of the provided text.
	Summarize(ctx context.Context, text string) (AnalysisResult, error)

	// Classify categorizes the text into predefined labels (e.g., "Work", "Newsletter", "Spam").
	Classify(ctx context.Context, text string) (string, error)

	// AnalyzeSentiment determines the sentiment and urgency of the text.
	// AnalyzeSentiment determines the sentiment and urgency of the text.
	AnalyzeSentiment(ctx context.Context, text string) (SentimentResult, error)

	// GenerateDraftReply generates a draft email reply based on the original email content and a user prompt.
	GenerateDraftReply(ctx context.Context, emailContent, userPrompt string) (string, error)
}

// EmbeddingProvider defines the interface for generating vector embeddings.
type EmbeddingProvider interface {
	// Embed generates a vector for a single text.
	Embed(ctx context.Context, text string) ([]float32, error)

	// EmbedBatch generates vectors for multiple texts.
	EmbedBatch(ctx context.Context, texts []string) ([][]float32, error)
}

type AnalysisResult struct {
	Summary     string   `json:"summary"`
	Category    string   `json:"category"`
	Sentiment   string   `json:"sentiment"`
	Urgency     string   `json:"urgency"`
	ActionItems []string `json:"action_items"`
}

type SentimentResult struct {
	Sentiment string // Positive, Neutral, Negative
	Urgency   string // High, Medium, Low
}
