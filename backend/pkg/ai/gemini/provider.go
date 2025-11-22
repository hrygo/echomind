package gemini

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"github.com/hrygo/echomind/pkg/ai"
	"google.golang.org/api/option"
)

type Provider struct {
	client         *genai.Client
	model          string
	embeddingModel string
	prompts        map[string]string
}

func NewProvider(ctx context.Context, settings map[string]interface{}, prompts map[string]string) (*Provider, error) {
	apiKey, _ := settings["api_key"].(string)
	modelName, _ := settings["model"].(string)
	embeddingModelName, _ := settings["embedding_model"].(string)

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}
	return &Provider{
		client:         client,
		model:          modelName,
		embeddingModel: embeddingModelName,
		prompts:        prompts,
	}, nil
}

func (p *Provider) Summarize(ctx context.Context, text string) (ai.AnalysisResult, error) {
	systemPrompt := p.prompts["summary"]
	if systemPrompt == "" {
		return ai.AnalysisResult{}, errors.New("summary prompt not configured")
	}

	model := p.client.GenerativeModel(p.model)
	model.ResponseMIMEType = "application/json"
	model.SystemInstruction = genai.NewUserContent(genai.Text(systemPrompt))

	resp, err := model.GenerateContent(ctx, genai.Text(text))
	if err != nil {
		return ai.AnalysisResult{}, err
	}

	response := extractText(resp)
	cleaned := cleanMarkdown(response)

	var result ai.AnalysisResult
	if err := json.Unmarshal([]byte(cleaned), &result); err != nil {
		return ai.AnalysisResult{Summary: response}, nil
	}

	return result, nil
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


func (p *Provider) GenerateDraftReply(ctx context.Context, emailContent, userPrompt string) (string, error) {
	systemPrompt := p.prompts["draft_reply"]
	if systemPrompt == "" {
		return "", errors.New("draft_reply prompt not configured")
	}

	fullPrompt := fmt.Sprintf("%s\n\nOriginal Email:\n%s\n\nUser Prompt:\n%s", systemPrompt, emailContent, userPrompt)
	return p.generateContent(ctx, fullPrompt, "")
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

// Embed generates a vector for a single text.
func (p *Provider) Embed(ctx context.Context, text string) ([]float32, error) {
	em := p.client.EmbeddingModel(p.embeddingModel)
	res, err := em.EmbedContent(ctx, genai.Text(text))
	if err != nil {
		return nil, err
	}
	return res.Embedding.Values, nil
}

// EmbedBatch generates vectors for multiple texts.
func (p *Provider) EmbedBatch(ctx context.Context, texts []string) ([][]float32, error) {
	em := p.client.EmbeddingModel(p.embeddingModel)
	batch := em.NewBatch()
	for _, text := range texts {
		batch.AddContent(genai.Text(text))
	}
	res, err := em.BatchEmbedContents(ctx, batch)
	if err != nil {
		return nil, err
	}
	var embeddings [][]float32
	for _, e := range res.Embeddings {
		embeddings = append(embeddings, e.Values)
	}
	return embeddings, nil
}
