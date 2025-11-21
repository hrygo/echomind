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

// NewAIProvider creates an AIProvider based on the configuration.
func NewAIProvider(cfg *configs.AIConfig) (ai.AIProvider, error) {
	prompts := toPromptMap(cfg.Prompts)

	switch cfg.Provider {
	case "openai":
		return openai.NewProvider(cfg.OpenAI.APIKey, cfg.OpenAI.Model, cfg.OpenAI.BaseURL, prompts), nil
	case "gemini":
		provider, err := gemini.NewProvider(context.Background(), cfg.Gemini.APIKey, cfg.Gemini.Model, prompts)
		if err != nil {
			return nil, err
		}
		return provider, nil
	case "deepseek":
		return deepseek.NewProvider(cfg.Deepseek.APIKey, cfg.Deepseek.Model, cfg.Deepseek.BaseURL, prompts), nil
	default:
		return nil, fmt.Errorf("unsupported AI provider: %s", cfg.Provider)
	}
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
