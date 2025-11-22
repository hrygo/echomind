package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hrygo/echomind/pkg/ai"
)

type ChatService struct {
	aiProvider    ai.AIProvider
	searchService *SearchService
}

func NewChatService(aiProvider ai.AIProvider, searchService *SearchService) *ChatService {
	return &ChatService{
		aiProvider:    aiProvider,
		searchService: searchService,
	}
}

func (s *ChatService) StreamChat(ctx context.Context, userID uuid.UUID, messages []ai.Message, ch chan<- ai.ChatCompletionChunk) error {
	// 1. Extract the last user message
	if len(messages) == 0 {
		return fmt.Errorf("no messages provided")
	}
	lastMsg := messages[len(messages)-1]
	if lastMsg.Role != "user" {
		return fmt.Errorf("last message must be from user")
	}

	// 2. RAG: Search for relevant context
	// We use a simple heuristic: search using the last user message
	// In a more advanced version, we might summarize the conversation history first.
	var searchResults []SearchResult
	var err error
	if s.searchService != nil {
		searchResults, err = s.searchService.Search(ctx, userID, lastMsg.Content, SearchFilters{}, 3)
		if err != nil {
			// Log error but continue
		}
	}

	// 3. Construct System Prompt with Context
	var contextBuilder strings.Builder
	if len(searchResults) > 0 {
		contextBuilder.WriteString("Here is some relevant information from the user's emails:\n\n")
		for _, result := range searchResults {
			contextBuilder.WriteString(fmt.Sprintf("Subject: %s\nFrom: %s\nDate: %s\nContent: %s\n\n",
				result.Subject, result.Sender, result.Date.Format("2006-01-02"), result.Snippet))
		}
		contextBuilder.WriteString("Answer the user's question based on the above context if relevant. If the context doesn't contain the answer, say so, but you can still answer general questions.\n")
	} else {
		contextBuilder.WriteString("You are EchoMind Copilot, a helpful AI assistant for managing emails and work.\n")
	}

	// 4. Prepare messages for AI Provider
	// We prepend the system prompt.
	// Note: Some providers (like OpenAI) accept a "system" role message at the beginning.
	// Gemini handles it via SystemInstruction, but our StreamChat implementation in Gemini
	// skips "system" messages in history.
	// So we should probably prepend the context to the *last* user message or the *first* user message?
	// Or we can just insert a system message at the start and let the provider handle it (OpenAI does, Gemini we might need to adjust).

	// Let's try prepending a system message.
	systemMsg := ai.Message{
		Role:    "system",
		Content: contextBuilder.String(),
	}

	finalMessages := append([]ai.Message{systemMsg}, messages...)

	// 5. Stream response
	return s.aiProvider.StreamChat(ctx, finalMessages, ch)
}
