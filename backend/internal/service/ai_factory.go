package service

import (
	"context"
	"fmt"

	"github.com/hrygo/echomind/configs"
	"github.com/hrygo/echomind/pkg/ai"
	"github.com/hrygo/echomind/pkg/ai/gemini"
	"github.com/hrygo/echomind/pkg/ai/mock"
	"github.com/hrygo/echomind/pkg/ai/openai"
)

// CompositeProvider combines an AIProvider (for text) and an EmbeddingProvider.
type CompositeProvider struct {
	ai.AIProvider
	ai.EmbeddingProvider
}

// NewAIProvider creates an AIProvider based on the configuration.
func NewAIProvider(cfg *configs.AIConfig) (ai.AIProvider, error) {
	prompts := toPromptMap(cfg.Prompts)
	var mainProvider ai.AIProvider
	var embeddingProvider ai.EmbeddingProvider
	var err error

	// Helper to create provider by name
	createProvider := func(name string) (interface{}, error) {
		// Handle "mock" special case
		if name == "mock" {
			return mock.NewProvider(), nil
		}

		pConfig, ok := cfg.Providers[name]
		if !ok {
			return nil, fmt.Errorf("provider configuration not found: %s", name)
		}

		switch pConfig.Protocol {
		case "openai":
			return openai.NewProvider(pConfig.Settings, prompts), nil
		case "gemini":
			return gemini.NewProvider(context.Background(), pConfig.Settings, prompts)
		default:
			return nil, fmt.Errorf("unsupported protocol: %s", pConfig.Protocol)
		}
	}

	// 1. Initialize Chat Provider
	chatProviderName := cfg.ActiveServices.Chat
	if chatProviderName == "" {
		chatProviderName = "mock" // Default fallback
	}
	
	chatP, err := createProvider(chatProviderName)
	if err != nil {
		return nil, fmt.Errorf("failed to create chat provider '%s': %w", chatProviderName, err)
	}
	
	if p, ok := chatP.(ai.AIProvider); ok {
		mainProvider = p
	} else {
		return nil, fmt.Errorf("provider '%s' does not implement AIProvider", chatProviderName)
	}

	// 2. Initialize Embedding Provider
	embedProviderName := cfg.ActiveServices.Embedding
	if embedProviderName == "" {
		embedProviderName = chatProviderName // Fallback to same provider
	}

	// Optimization: If chat and embedding use same provider name, reuse instance if it implements both
	if embedProviderName == chatProviderName {
		if p, ok := chatP.(ai.EmbeddingProvider); ok {
			embeddingProvider = p
		} else {
			// Configured same provider but it doesn't support embeddings? 
			// Try re-instantiating (maybe it needs different settings? unlikely but safe)
			// Actually, let's just error or fallback to mock? 
			// Better to try creating it again or logging warning. 
			// For now, let's assume if it's the same name, it SHOULD support it, or we try creating it.
		}
	}
	
	if embeddingProvider == nil {
		embedP, err := createProvider(embedProviderName)
		if err != nil {
			return nil, fmt.Errorf("failed to create embedding provider '%s': %w", embedProviderName, err)
		}
		if p, ok := embedP.(ai.EmbeddingProvider); ok {
			embeddingProvider = p
		} else {
			return nil, fmt.Errorf("provider '%s' does not implement EmbeddingProvider", embedProviderName)
		}
	}

	return &CompositeProvider{
		AIProvider:        mainProvider,
		EmbeddingProvider: embeddingProvider,
	}, nil
}

// toPromptMap converts a PromptConfig struct to a map[string]string.
func toPromptMap(pc configs.PromptConfig) map[string]string {
	return map[string]string{
		"summary":    pc.Summary,
		"classify":   pc.Classify,
		"sentiment":  pc.Sentiment,
		"draft_reply": pc.DraftReply,
	}
}
