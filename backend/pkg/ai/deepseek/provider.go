package deepseek

import (
	"echomind.com/backend/pkg/ai/openai"
)

// DeepSeek is OpenAI Compatible, so we can reuse the OpenAI provider with a custom BaseURL.
// This reduces code duplication significantly.

func NewProvider(apiKey, model, baseURL string, prompts map[string]string) *openai.Provider {
    if baseURL == "" {
        baseURL = "https://api.deepseek.com/v1"
    }
	return openai.NewProvider(apiKey, model, baseURL, prompts)
}
