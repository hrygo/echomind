package service

import (
	"context"
	"fmt"

	"github.com/hrygo/echomind/pkg/ai"
	"github.com/hrygo/echomind/pkg/ai/deepseek"
	"github.com/hrygo/echomind/pkg/ai/gemini"
	"github.com/hrygo/echomind/pkg/ai/openai"
	"github.com/spf13/viper"
)

type SummaryService struct {
	provider ai.AIProvider
}

func NewSummaryService(provider ai.AIProvider) *SummaryService {
	return &SummaryService{
		provider: provider,
	}
}

func (s *SummaryService) GenerateSummary(ctx context.Context, text string) (string, error) {
	return s.provider.Summarize(ctx, text)
}

func (s *SummaryService) AnalyzeSentiment(ctx context.Context, text string) (ai.SentimentResult, error) {
	return s.provider.AnalyzeSentiment(ctx, text)
}

// AIProviderFactory creates an AI provider based on configuration.
func AIProviderFactory(v *viper.Viper) (ai.AIProvider, error) {
	providerType := v.GetString("ai.provider")
	prompts := v.GetStringMapString("ai.prompts")
	
	switch providerType {
	case "deepseek":
		apiKey := v.GetString("ai.deepseek.api_key")
		model := v.GetString("ai.deepseek.model")
		baseURL := v.GetString("ai.deepseek.base_url")
		return deepseek.NewProvider(apiKey, model, baseURL, prompts), nil
	case "openai":
		apiKey := v.GetString("ai.openai.api_key")
		model := v.GetString("ai.openai.model")
		baseURL := v.GetString("ai.openai.base_url")
		return openai.NewProvider(apiKey, model, baseURL, prompts), nil
	case "gemini":
		apiKey := v.GetString("ai.gemini.api_key")
		model := v.GetString("ai.gemini.model")
		// Gemini provider needs context for initialization
		return gemini.NewProvider(context.Background(), apiKey, model, prompts)
	default:
		return nil, fmt.Errorf("unsupported AI provider: %s", providerType)
	}
}
