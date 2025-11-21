package gemini

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/hrygo/echomind/pkg/ai"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type Provider struct {
	client  *genai.Client
	model   string
	prompts map[string]string
}

func NewProvider(ctx context.Context, apiKey, modelName string, prompts map[string]string) (*Provider, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}
	return &Provider{
		client:  client,
		model:   modelName,
		prompts: prompts,
	}, nil
}

func (p *Provider) Summarize(ctx context.Context, text string) (string, error) {
	systemPrompt := p.prompts["summary"]
	if systemPrompt == "" {
		return "", errors.New("summary prompt not configured")
	}
	return p.generateContent(ctx, systemPrompt, text)
}

func (p *Provider) Classify(ctx context.Context, text string) (string, error) {
	systemPrompt := p.prompts["classify"]
	if systemPrompt == "" {
		return "", errors.New("classify prompt not configured")
	}
	return p.generateContent(ctx, systemPrompt, text)
}

func (p *Provider) AnalyzeSentiment(ctx context.Context, text string) (ai.SentimentResult, error) {
	systemPrompt := p.prompts["sentiment"]
	if systemPrompt == "" {
		return ai.SentimentResult{}, errors.New("sentiment prompt not configured")
	}

	model := p.client.GenerativeModel(p.model)
    model.ResponseMIMEType = "application/json" // Force JSON mode for Gemini
    
    prompt := fmt.Sprintf("%s\n\nEmail Content:\n%s", systemPrompt, text)
    resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return ai.SentimentResult{}, err
	}
    
    response := extractText(resp)

	var result struct {
		Sentiment string `json:"sentiment"`
		Urgency   string `json:"urgency"`
	}
    
    // Gemini usually respects JSON mode well, but we can clean just in case
    cleaned := cleanMarkdown(response)

	if err := json.Unmarshal([]byte(cleaned), &result); err != nil {
		return ai.SentimentResult{Sentiment: "Neutral", Urgency: "Medium"}, nil
	}

	return ai.SentimentResult{
		Sentiment: result.Sentiment,
		Urgency:   result.Urgency,
	}, nil
}

func (p *Provider) generateContent(ctx context.Context, systemPrompt, userContent string) (string, error) {
	model := p.client.GenerativeModel(p.model)
    model.SystemInstruction = genai.NewUserContent(genai.Text(systemPrompt))
    
	resp, err := model.GenerateContent(ctx, genai.Text(userContent))
	if err != nil {
		return "", err
	}

	return extractText(resp), nil
}

func extractText(resp *genai.GenerateContentResponse) string {
	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return ""
	}
	
    var sb strings.Builder
    for _, part := range resp.Candidates[0].Content.Parts {
        if txt, ok := part.(genai.Text); ok {
            sb.WriteString(string(txt))
        }
    }
    return sb.String()
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