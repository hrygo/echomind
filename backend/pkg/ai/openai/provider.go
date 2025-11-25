package openai

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/hrygo/echomind/configs"
	"github.com/hrygo/echomind/pkg/ai"
	"github.com/hrygo/echomind/pkg/ai/registry"
	openai "github.com/sashabaranov/go-openai"
)

func init() {
	registry.Register("openai", NewProvider)
}

type Provider struct {
	client         *openai.Client
	model          string
	embeddingModel string
	dimensions     int
	prompts        map[string]string
}

func NewProvider(ctx context.Context, settings configs.ProviderSettings, prompts map[string]string) (ai.AIProvider, error) {
	apiKey, _ := settings["api_key"].(string)
	model, _ := settings["model"].(string)
	baseURL, _ := settings["base_url"].(string)
	embeddingModel, _ := settings["embedding_model"].(string)
	dimensions := 1024 // default dimension

	// Get embedding dimensions from config if specified
	if dim, ok := settings["embedding_dimensions"].(int); ok {
		dimensions = dim
	} else if dimFloat, ok := settings["embedding_dimensions"].(float64); ok {
		dimensions = int(dimFloat)
	}

	// Default embedding model if not specified
	if embeddingModel == "" {
		embeddingModel = string(openai.SmallEmbedding3)
	}

	config := openai.DefaultConfig(apiKey)
	if baseURL != "" {
		config.BaseURL = baseURL
	}
	return &Provider{
		client:         openai.NewClientWithConfig(config),
		model:          model,
		embeddingModel: embeddingModel,
		dimensions:     dimensions,
		prompts:        prompts,
	}, nil
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

func (p *Provider) GenerateDraftReply(ctx context.Context, emailContent, userPrompt string) (string, error) {
	systemPrompt := p.prompts["draft_reply"]
	if systemPrompt == "" {
		systemPrompt = "You are an email assistant. Generate a professional email reply based on the provided email content and user instructions."
	}

	// Ensure we have valid content
	if emailContent == "" {
		emailContent = "No email content provided."
	}
	if userPrompt == "" {
		userPrompt = "Generate a brief, professional email reply."
	}

	fullUserPrompt := fmt.Sprintf("Original Email:\n%s\n\nUser Instructions:\n%s", emailContent, userPrompt)
	return p.chatCompletion(ctx, systemPrompt, fullUserPrompt, false)
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

func (p *Provider) StreamChat(ctx context.Context, messages []ai.Message, ch chan<- ai.ChatCompletionChunk) error {
	defer close(ch)

	var openaiMessages []openai.ChatCompletionMessage
	for _, msg := range messages {
		openaiMessages = append(openaiMessages, openai.ChatCompletionMessage{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	req := openai.ChatCompletionRequest{
		Model:    p.model,
		Messages: openaiMessages,
		Stream:   true,
	}

	stream, err := p.client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		return err
	}
	defer stream.Close()

	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return err
		}

		if len(response.Choices) > 0 {
			ch <- ai.ChatCompletionChunk{
				ID: response.ID,
				Choices: []ai.Choice{
					{
						Index: response.Choices[0].Index,
						Delta: ai.DeltaContent{
							Content: response.Choices[0].Delta.Content,
						},
					},
				},
			}
		}
	}
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

// Embed generates a vector for a single text using the configured embedding model.
func (p *Provider) Embed(ctx context.Context, text string) ([]float32, error) {
	req := openai.EmbeddingRequest{
		Input: []string{text},
		Model: openai.EmbeddingModel(p.embeddingModel),
	}

	resp, err := p.client.CreateEmbeddings(ctx, req)
	if err != nil {
		return nil, err
	}

	if len(resp.Data) == 0 {
		return nil, errors.New("no embedding data returned")
	}

	return resp.Data[0].Embedding, nil
}

// EmbedBatch generates vectors for multiple texts.
func (p *Provider) EmbedBatch(ctx context.Context, texts []string) ([][]float32, error) {
	req := openai.EmbeddingRequest{
		Input: texts,
		Model: openai.EmbeddingModel(p.embeddingModel),
	}

	resp, err := p.client.CreateEmbeddings(ctx, req)
	if err != nil {
		return nil, err
	}

	var embeddings [][]float32
	for _, data := range resp.Data {
		embeddings = append(embeddings, data.Embedding)
	}

	return embeddings, nil
}

// GetDimensions returns the dimension size of the vectors generated by this provider.
func (p *Provider) GetDimensions() int {
	return p.dimensions
}
