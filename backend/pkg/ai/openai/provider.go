package openai

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/hrygo/echomind/pkg/ai"
	openai "github.com/sashabaranov/go-openai"
)

type Provider struct {
	client  *openai.Client
	model   string
	prompts map[string]string
}

func NewProvider(apiKey, model, baseURL string, prompts map[string]string) *Provider {
	config := openai.DefaultConfig(apiKey)
	if baseURL != "" {
		config.BaseURL = baseURL
	}
	return &Provider{
		client:  openai.NewClientWithConfig(config),
		model:   model,
		prompts: prompts,
	}
}

func (p *Provider) Summarize(ctx context.Context, text string) (ai.AnalysisResult, error) {
	systemPrompt := p.prompts["summary"]
	if systemPrompt == "" {
		return ai.AnalysisResult{}, errors.New("summary prompt not configured")
	}

	response, err := p.chatCompletion(ctx, systemPrompt, text, true)
	if err != nil {
		return ai.AnalysisResult{}, err
	}

	cleaned := cleanMarkdown(response)
	var result ai.AnalysisResult
	if err := json.Unmarshal([]byte(cleaned), &result); err != nil {
		// Fallback if JSON parsing fails
		return ai.AnalysisResult{Summary: response}, nil
	}

	return result, nil
}

func (p *Provider) Classify(ctx context.Context, text string) (string, error) {
	systemPrompt := p.prompts["classify"]
	if systemPrompt == "" {
		return "", errors.New("classify prompt not configured")
	}
	return p.chatCompletion(ctx, systemPrompt, text, false)
}

func (p *Provider) AnalyzeSentiment(ctx context.Context, text string) (ai.SentimentResult, error) {
	systemPrompt := p.prompts["sentiment"]
	if systemPrompt == "" {
		return ai.SentimentResult{}, errors.New("sentiment prompt not configured")
	}

	response, err := p.chatCompletion(ctx, systemPrompt, text, true)
	if err != nil {
		return ai.SentimentResult{}, err
	}

	var result struct {
		Sentiment string `json:"sentiment"`
		Urgency   string `json:"urgency"`
	}

	// Clean markdown if present
	cleaned := cleanMarkdown(response)

	if err := json.Unmarshal([]byte(cleaned), &result); err != nil {
		return ai.SentimentResult{Sentiment: "Neutral", Urgency: "Medium"}, nil
	}

	return ai.SentimentResult{
		Sentiment: result.Sentiment,
		Urgency:   result.Urgency,
	}, nil
}

func (p *Provider) chatCompletion(ctx context.Context, systemPrompt, userContent string, jsonMode bool) (string, error) {
	req := openai.ChatCompletionRequest{
		Model: p.model,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: systemPrompt},
			{Role: openai.ChatMessageRoleUser, Content: userContent},
		},
	}

	if jsonMode {
		req.ResponseFormat = &openai.ChatCompletionResponseFormat{
			Type: openai.ChatCompletionResponseFormatTypeJSONObject,
		}
	}

	resp, err := p.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", errors.New("no choices returned from OpenAI API")
	}

	return resp.Choices[0].Message.Content, nil
}

func cleanMarkdown(text string) string {
	cleaned := text
	if len(cleaned) > 7 && cleaned[:7] == "```json" {
		cleaned = cleaned[7:]
		if len(cleaned) > 3 && cleaned[len(cleaned)-3:] == "```" {
			cleaned = cleaned[:len(cleaned)-3]
		}
	} else if len(cleaned) > 3 && cleaned[:3] == "```" {
		cleaned = cleaned[3:]
		if len(cleaned) > 3 && cleaned[len(cleaned)-3:] == "```" {
			cleaned = cleaned[:len(cleaned)-3]
		}
	}
	return strings.TrimSpace(cleaned)
}
