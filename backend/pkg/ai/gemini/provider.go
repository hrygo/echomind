package gemini

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"github.com/hrygo/echomind/configs"
	"github.com/hrygo/echomind/pkg/ai"
	"github.com/hrygo/echomind/pkg/ai/registry"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

func init() {
	registry.Register("gemini", NewProvider)
}

type Provider struct {
	client         *genai.Client
	model          string
	embeddingModel string
	prompts        map[string]string
}

func NewProvider(ctx context.Context, settings configs.ProviderSettings, prompts map[string]string) (ai.AIProvider, error) {
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
		systemPrompt = "You are an email assistant. Generate a professional email reply based on the provided email content and user instructions."
	}

	// Ensure we have valid content
	if emailContent == "" {
		emailContent = "No email content provided."
	}
	if userPrompt == "" {
		userPrompt = "Generate a brief, professional email reply."
	}

	fullPrompt := fmt.Sprintf("%s\n\nOriginal Email:\n%s\n\nUser Instructions:\n%s", systemPrompt, emailContent, userPrompt)
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

func (p *Provider) StreamChat(ctx context.Context, messages []ai.Message, ch chan<- ai.ChatCompletionChunk) error {
	defer close(ch)

	model := p.client.GenerativeModel(p.model)
	cs := model.StartChat()

	// Convert history (excluding the last message which is the new prompt)
	if len(messages) > 1 {
		var history []*genai.Content
		for _, msg := range messages[:len(messages)-1] {
			role := "user"
			if msg.Role == "assistant" {
				role = "model"
			} else if msg.Role == "system" {
				// Gemini doesn't support system messages in history directly in the same way,
				// usually set as SystemInstruction on the model.
				// For simplicity here, we might skip or prepend to next user message.
				// But since we set SystemInstruction in other methods, let's assume system prompt is handled via config if needed.
				// Here we just handle user/model turn.
				continue
			}
			history = append(history, &genai.Content{
				Parts: []genai.Part{genai.Text(msg.Content)},
				Role:  role,
			})
		}
		cs.History = history
	}

	lastMsg := messages[len(messages)-1]
	iter := cs.SendMessageStream(ctx, genai.Text(lastMsg.Content))

	for i := 0; ; i++ {
		resp, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			return nil
		}
		if err != nil {
			return err
		}

		// Use the new function to extract both content and widget data
		content, widgetData := extractContentAndWidget(resp)

		if content != "" || widgetData != nil { // Send if either content or widget is present
			delta := ai.DeltaContent{}
			if content != "" {
				delta.Content = content
			}
			if widgetData != nil {
				delta.Widget = widgetData
			}

			chunk := ai.ChatCompletionChunk{
				ID: fmt.Sprintf("chatcmpl-%d", i), // Simple ID, can be UUID
				Choices: []ai.Choice{
					{
						Index: 0,
						Delta: delta,
					},
				},
			}
			ch <- chunk
		}
	}
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

// extractContentAndWidget parses the AI response for potential widget data.
// It returns the remaining text content and an optional WidgetData object.
func extractContentAndWidget(resp *genai.GenerateContentResponse) (string, *ai.WidgetData) {
	fullText := extractText(resp)

	// Check for widget pattern: ```widget_TYPE\n{JSON_DATA}\n```
	widgetPattern := regexp.MustCompile("```widget_([a-zA-Z0-9_]+)\\n([\\s\\S]*?)\\n```")
	matches := widgetPattern.FindStringSubmatch(fullText)

	if len(matches) == 3 {
		widgetType := matches[1]
		jsonData := strings.TrimSpace(matches[2])

		var data map[string]interface{}
		if err := json.Unmarshal([]byte(jsonData), &data); err == nil {
			// Successfully parsed widget JSON
			widget := &ai.WidgetData{
				Type: widgetType,
				Data: data,
			}
			// Remove the widget block from the fullText
			cleanText := strings.Replace(fullText, matches[0], "", 1)
			return strings.TrimSpace(cleanText), widget
		}
	}
	// No widget found or failed to parse, return full text as content
	return fullText, nil
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
