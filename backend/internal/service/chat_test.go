package service

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hrygo/echomind/internal/model"
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

// MockContextSearcher mocks ContextSearcher interface
type MockContextSearcher struct {
	mock.Mock
}

func (m *MockContextSearcher) Search(ctx context.Context, userID uuid.UUID, query string, filters SearchFilters, limit int) ([]SearchResult, error) {
	args := m.Called(ctx, userID, query, filters, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]SearchResult), args.Error(1)
}

// MockEmailRetriever mocks EmailRetriever interface
type MockEmailRetriever struct {
	mock.Mock
}

func (m *MockEmailRetriever) GetEmailsByIDs(ctx context.Context, userID uuid.UUID, emailIDs []uuid.UUID) ([]model.Email, error) {
	args := m.Called(ctx, userID, emailIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Email), args.Error(1)
}

func TestChatService_StreamChat(t *testing.T) {
	// Setup Mocks
	mockAI := new(MockAIProvider)
	mockSearch := new(MockContextSearcher)
	mockEmail := new(MockEmailRetriever)

	chatService := NewChatService(mockAI, mockSearch, mockEmail)

	userID := uuid.New()
	ctx := context.Background()

	t.Run("Basic Chat without Context", func(t *testing.T) {
		messages := []ai.Message{{Role: "user", Content: "Hello"}}
		ch := make(chan ai.ChatCompletionChunk, 10)

		// Expect Search to be called (Auto-Search Fallback) but return empty
		mockSearch.On("Search", ctx, userID, "Hello", mock.Anything, 3).Return([]SearchResult{}, nil).Once()

		// Expect StreamChat to be called
		mockAI.On("StreamChat", ctx, mock.MatchedBy(func(msgs []ai.Message) bool {
			// Check if system prompt is generic (no context found)
			return len(msgs) == 2 && msgs[0].Role == "system" && msgs[1].Content == "Hello"
		}), (chan<- ai.ChatCompletionChunk)(ch)).Return(nil).Run(func(args mock.Arguments) {
			chArg := args.Get(2).(chan<- ai.ChatCompletionChunk)
			chArg <- ai.ChatCompletionChunk{ID: "1", Choices: []ai.Choice{{Index: 0, Delta: ai.DeltaContent{Content: "Hi"}}}}
			close(chArg)
		}).Once()

		err := chatService.StreamChat(ctx, userID, messages, nil, ch)
		assert.NoError(t, err)

		// Verify mocks
		mockSearch.AssertExpectations(t)
		mockAI.AssertExpectations(t)
	})

	t.Run("Chat with Explicit ContextRefIDs", func(t *testing.T) {
		emailID := uuid.New()
		contextRefIDs := []uuid.UUID{emailID}
		messages := []ai.Message{{Role: "user", Content: "Summarize this email"}}
		ch := make(chan ai.ChatCompletionChunk, 10)

		// Mock Email Retrieval
		mockEmail.On("GetEmailsByIDs", ctx, userID, contextRefIDs).Return([]model.Email{
			{
				ID:      emailID,
				Subject: "Test Email",
				Sender:  "sender@example.com",
				Snippet: "This is a test email content.",
				Date:    time.Now(),
			},
		}, nil).Once()

		// Expect StreamChat to be called with context injected
		mockAI.On("StreamChat", ctx, mock.MatchedBy(func(msgs []ai.Message) bool {
			// Check if system prompt contains the email content
			if len(msgs) < 2 {
				return false
			}
			systemPrompt := msgs[0].Content
			return msgs[0].Role == "system" &&
				strings.Contains(systemPrompt, "Test Email") &&
				strings.Contains(systemPrompt, "This is a test email content")
		}), (chan<- ai.ChatCompletionChunk)(ch)).Return(nil).Run(func(args mock.Arguments) {
			chArg := args.Get(2).(chan<- ai.ChatCompletionChunk)
			chArg <- ai.ChatCompletionChunk{ID: "1", Choices: []ai.Choice{{Index: 0, Delta: ai.DeltaContent{Content: "Summary..."}}}}
			close(chArg)
		}).Once()

		err := chatService.StreamChat(ctx, userID, messages, contextRefIDs, ch)
		assert.NoError(t, err)

		mockEmail.AssertExpectations(t)
		mockAI.AssertExpectations(t)
	})

	t.Run("Chat with Auto-Search Fallback", func(t *testing.T) {
		messages := []ai.Message{{Role: "user", Content: "What about the budget?"}}
		ch := make(chan ai.ChatCompletionChunk, 10)

		// Mock Search returning results
		mockSearch.On("Search", ctx, userID, "What about the budget?", mock.Anything, 3).Return([]SearchResult{
			{
				Subject: "Budget Report",
				Snippet: "The budget is tight.",
				Sender:  "finance@example.com",
				Date:    time.Now(),
			},
		}, nil).Once()

		// Expect StreamChat with injected context
		mockAI.On("StreamChat", ctx, mock.MatchedBy(func(msgs []ai.Message) bool {
			systemPrompt := msgs[0].Content
			return strings.Contains(systemPrompt, "Budget Report") && strings.Contains(systemPrompt, "The budget is tight")
		}), (chan<- ai.ChatCompletionChunk)(ch)).Return(nil).Run(func(args mock.Arguments) {
			chArg := args.Get(2).(chan<- ai.ChatCompletionChunk)
			close(chArg)
		}).Once()

		err := chatService.StreamChat(ctx, userID, messages, nil, ch)
		assert.NoError(t, err)

		mockSearch.AssertExpectations(t)
		mockAI.AssertExpectations(t)
	})
}
