package service

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/hrygo/echomind/pkg/ai"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAIProvider is a mock implementation of ai.AIProvider
type MockAIProvider struct {
	mock.Mock
}

func (m *MockAIProvider) Summarize(ctx context.Context, text string) (ai.AnalysisResult, error) {
	args := m.Called(ctx, text)
	return args.Get(0).(ai.AnalysisResult), args.Error(1)
}

func (m *MockAIProvider) Classify(ctx context.Context, text string) (string, error) {
	args := m.Called(ctx, text)
	return args.String(0), args.Error(1)
}

func (m *MockAIProvider) AnalyzeSentiment(ctx context.Context, text string) (ai.SentimentResult, error) {
	args := m.Called(ctx, text)
	return args.Get(0).(ai.SentimentResult), args.Error(1)
}

func (m *MockAIProvider) GenerateDraftReply(ctx context.Context, emailContent, userPrompt string) (string, error) {
	args := m.Called(ctx, emailContent, userPrompt)
	return args.String(0), args.Error(1)
}

func (m *MockAIProvider) StreamChat(ctx context.Context, messages []ai.Message, ch chan<- string) error {
	args := m.Called(ctx, messages, ch)
	return args.Error(0)
}

func TestChatService_StreamChat(t *testing.T) {
	mockAI := new(MockAIProvider)
	// We can mock SearchService too, but for simplicity let's pass nil or a mock if needed.
	// Since SearchService struct doesn't implement an interface in the current code (it's a struct),
	// we might need to refactor SearchService to an interface or just pass nil if we don't test RAG here.
	// However, ChatService uses *SearchService.
	// For this test, let's assume SearchService can be nil or we just test the flow where search fails or returns empty.
	// But ChatService calls s.searchService.Search. If s.searchService is nil, it will panic.
	// So we need a real SearchService or refactor to interface.
	// Given the constraints, let's skip RAG part testing or mock the SearchService if possible.
	// Since we can't easily mock a struct method without an interface, let's create a ChatService with a nil SearchService
	// AND modify ChatService to handle nil SearchService gracefully, OR we just test the AI part.

	// Actually, let's modify ChatService to check for nil SearchService to make it testable without DB.
	chatService := NewChatService(mockAI, nil)

	ctx := context.Background()
	userID := uuid.New()
	messages := []ai.Message{{Role: "user", Content: "Hello"}}
	ch := make(chan string, 10)

	// Expect StreamChat to be called
	mockAI.On("StreamChat", mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		// Get the channel from arguments
		chArg := args.Get(2).(chan<- string)
		chArg <- "Hello"
		chArg <- " World"
		close(chArg)
	})

	// We need to modify ChatService to handle nil SearchService first for this test to pass without panic.
	// Or we can try to instantiate a SearchService with a mock DB? That's too complex for now.
	// Let's modify ChatService to be robust.

	err := chatService.StreamChat(ctx, userID, messages, ch)
	assert.NoError(t, err)

	// Collect output
	var output string
	for msg := range ch {
		output += msg
	}
	assert.Equal(t, "Hello World", output)
	mockAI.AssertExpectations(t)
}
