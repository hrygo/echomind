package service

import (
	"context"
	"fmt"

	"github.com/hrygo/echomind/configs"
	"github.com/hrygo/echomind/pkg/ai"
	"github.com/hrygo/echomind/pkg/ai/deepseek"
	"github.com/hrygo/echomind/pkg/ai/gemini"
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
	var err error

	// 1. Initialize Main Provider (Chat/Summary)
	switch cfg.Provider {
	case "openai":
		mainProvider = openai.NewProvider(cfg.OpenAI.APIKey, cfg.OpenAI.Model, cfg.OpenAI.BaseURL, prompts)
	case "gemini":
		mainProvider, err = gemini.NewProvider(context.Background(), cfg.Gemini.APIKey, cfg.Gemini.Model, prompts)
		if err != nil {
			return nil, err
		}
	case "deepseek":
		mainProvider = deepseek.NewProvider(cfg.Deepseek.APIKey, cfg.Deepseek.Model, cfg.Deepseek.BaseURL, prompts)
	default:
		return nil, fmt.Errorf("unsupported AI provider: %s", cfg.Provider)
	}

	// 2. Initialize Embedding Provider
	// Default to main provider if not specified, unless main is DeepSeek (fallback to OpenAI for embeddings)
	embeddingProviderName := cfg.EmbeddingProvider
	if embeddingProviderName == "" {
		if cfg.Provider == "deepseek" {
			embeddingProviderName = "openai" // DeepSeek doesn't support embeddings yet, force OpenAI
		} else {
			embeddingProviderName = cfg.Provider
		}
	}

	var embedder ai.EmbeddingProvider
	
	// Determine embedding model (fallback to default if not specified)
	embeddingModel := cfg.EmbeddingModel
	if embeddingModel == "" {
		if embeddingProviderName == "openai" {
			embeddingModel = "text-embedding-3-small"
		} else if embeddingProviderName == "gemini" {
			embeddingModel = "text-embedding-004" // Common gemini embedding model
		}
	}

	switch embeddingProviderName {
	case "openai":
		// Re-use OpenAI config for embeddings, but allow overriding model
		embedder = openai.NewProvider(cfg.OpenAI.APIKey, embeddingModel, cfg.OpenAI.BaseURL, nil)
	case "gemini":
		// Re-use Gemini config
		gp, err := gemini.NewProvider(context.Background(), cfg.Gemini.APIKey, embeddingModel, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create gemini embedder: %w", err)
		}
		var asInterface interface{} = gp
		var ok bool
		embedder, ok = asInterface.(ai.EmbeddingProvider)
		if !ok {
			return nil, fmt.Errorf("gemini provider does not support embeddings")
		}
	case "deepseek":
		// If user explicitly requests DeepSeek embeddings
		embedder = deepseek.NewProvider(cfg.Deepseek.APIKey, embeddingModel, cfg.Deepseek.BaseURL, nil)
	default:
		// If main provider implements embedding, use it (fallback logic)
		if e, ok := mainProvider.(ai.EmbeddingProvider); ok {
			embedder = e
		} else {
			// Fallback to OpenAI if not specified and main doesn't support
			embedder = openai.NewProvider(cfg.OpenAI.APIKey, "text-embedding-3-small", cfg.OpenAI.BaseURL, nil)
		}
	}

	return &CompositeProvider{
		AIProvider:        mainProvider,
		EmbeddingProvider: embedder,
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
