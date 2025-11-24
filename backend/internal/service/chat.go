package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hrygo/echomind/internal/model"
	"github.com/hrygo/echomind/pkg/ai"
)

// ContextSearcher defines the interface for searching context (emails, etc).
type ContextSearcher interface {
	Search(ctx context.Context, userID uuid.UUID, query string, filters SearchFilters, limit int) ([]SearchResult, error)
}

// EmailRetriever defines the interface for retrieving emails by ID.
type EmailRetriever interface {
	GetEmailsByIDs(ctx context.Context, userID uuid.UUID, emailIDs []uuid.UUID) ([]model.Email, error)
}

type ChatService struct {
	aiProvider    ai.AIProvider
	searchService ContextSearcher
	emailService  EmailRetriever
}

func NewChatService(aiProvider ai.AIProvider, searchService ContextSearcher, emailService EmailRetriever) *ChatService {
	return &ChatService{
		aiProvider:    aiProvider,
		searchService: searchService,
		emailService:  emailService,
	}
}

func (s *ChatService) StreamChat(ctx context.Context, userID uuid.UUID, messages []ai.Message, contextRefIDs []uuid.UUID, ch chan<- ai.ChatCompletionChunk) error {
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
	var contextBuilder strings.Builder

	// Strategy A: Explicit Context (High Priority)
	if len(contextRefIDs) > 0 && s.emailService != nil {
		emails, err := s.emailService.GetEmailsByIDs(ctx, userID, contextRefIDs)
		if err == nil && len(emails) > 0 {
			contextBuilder.WriteString("The user has provided the following specific emails as context:\n\n")
			for _, email := range emails {
				// Use Snippet for now to save tokens, or truncated BodyText
				content := email.Snippet
				if content == "" {
					content = "(No Content)"
				}
				contextBuilder.WriteString(fmt.Sprintf("Subject: %s\nFrom: %s\nDate: %s\nContent: %s\n\n",
					email.Subject, email.Sender, email.Date.Format("2006-01-02"), content))
			}
		}
	} else if s.searchService != nil {
		// Strategy B: Auto-Search (Fallback)
		searchResults, err := s.searchService.Search(ctx, userID, lastMsg.Content, SearchFilters{}, 3)
		if err != nil {
			// Log error but continue
			_ = err
		} else if len(searchResults) > 0 {
			contextBuilder.WriteString("Here is some relevant information from the user's emails:\n\n")
			for _, result := range searchResults {
				contextBuilder.WriteString(fmt.Sprintf("Subject: %s\nFrom: %s\nDate: %s\nContent: %s\n\n",
					result.Subject, result.Sender, result.Date.Format("2006-01-02"), result.Snippet))
			}
		}
	}

	// 3. Construct System Prompt with Context
	if contextBuilder.Len() > 0 {
		contextBuilder.WriteString("Answer the user's question based on the above context if relevant. If the context doesn't contain the answer, say so, but you can still answer general questions.\n")
	} else {
		contextBuilder.WriteString("You are EchoMind Copilot, a helpful AI assistant for managing emails and work.\n")
	}

	// Add Widget Instructions
	contextBuilder.WriteString(`
You can generate interactive widgets for the user.
Supported widgets:
1. Task List: <widget type="task_list">[{"title": "Task 1", "due": "2023-10-27"}, ...]</widget>
2. Email Draft: <widget type="email_draft">{"to": "...", "subject": "...", "body": "..."}</widget>
3. Calendar Event: <widget type="calendar_event">{"title": "...", "start": "...", "end": "..."}</widget>

When the user asks to create tasks, draft emails, or schedule meetings, output the corresponding widget XML block.
`)

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
