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

func (m *MockAIProvider) StreamChat(ctx context.Context, messages []ai.Message, ch chan<- ai.ChatCompletionChunk) error {
	args := m.Called(ctx, messages, ch)
	return args.Error(0)
}

func TestChatService_StreamChat(t *testing.T) {
	mockAI := new(MockAIProvider)
	chatService := NewChatService(mockAI, nil) // SearchService can be nil for this test

	ctx := context.Background()
	userID := uuid.New()
	messages := []ai.Message{{Role: "user", Content: "Hello"}}
	ch := make(chan ai.ChatCompletionChunk, 10)

	// Expect StreamChat to be called
	mockAI.On("StreamChat", mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		// Get the channel from arguments
		chArg := args.Get(2).(chan<- ai.ChatCompletionChunk)
		// Send mock chunks
		chArg <- ai.ChatCompletionChunk{ID: "1", Choices: []ai.Choice{{Index: 0, Delta: ai.DeltaContent{Content: "Hello"}}}}
		chArg <- ai.ChatCompletionChunk{ID: "2", Choices: []ai.Choice{{Index: 0, Delta: ai.DeltaContent{Content: " World"}}}}
		close(chArg)
	})

	err := chatService.StreamChat(ctx, userID, messages, ch)
	assert.NoError(t, err)

	// Collect output
	var output string
	for chunk := range ch {
		if len(chunk.Choices) > 0 {
			output += chunk.Choices[0].Delta.Content
		}
	}
	assert.Equal(t, "Hello World", output)
	mockAI.AssertExpectations(t)
}
