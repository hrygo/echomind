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

	// StreamChat streams the chat response token by token as structured chunks.
	StreamChat(ctx context.Context, messages []Message, ch chan<- ChatCompletionChunk) error
}

// EmbeddingProvider defines the interface for generating vector embeddings.
type EmbeddingProvider interface {
	// Embed generates a vector for a single text.
	Embed(ctx context.Context, text string) ([]float32, error)

	// EmbedBatch generates vectors for multiple texts.
	EmbedBatch(ctx context.Context, texts []string) ([][]float32, error)
}

type Message struct {
	Role    string `json:"role"` // "system", "user", "assistant"
	Content string `json:"content"`
}

// ChatCompletionChunk represents a chunk of the streamed chat completion response.
type ChatCompletionChunk struct {
	ID      string  `json:"id"`
	Choices []Choice `json:"choices"`
}

type DeltaContent struct {
	Content string     `json:"content,omitempty"`
	Widget  *WidgetData `json:"widget,omitempty"`
}

type WidgetData struct {
	Type  string                 `json:"type"`  // e.g., "task_card", "search_result_card"
	Data  map[string]interface{} `json:"data"` // Arbitrary data for the widget
}

type Choice struct {
	Index int `json:"index"`
	Delta DeltaContent `json:"delta"`
}

type AnalysisResult struct {
	Summary      string        `json:"summary"`
	Category     string        `json:"category"`
	Sentiment    string        `json:"sentiment"`
	Urgency      string        `json:"urgency"`
	ActionItems  []string      `json:"action_items"`
	SmartActions []SmartAction `json:"smart_actions"`
}

type SmartAction struct {
	Type  string            `json:"type"`  // "calendar_event", "create_task"
	Label string            `json:"label"` // Display text for the button
	Data  map[string]string `json:"data"`  // Context data (title, date, etc.)
}

type SentimentResult struct {
	Sentiment string // Positive, Neutral, Negative
	Urgency   string // High, Medium, Low
}
