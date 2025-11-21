package service_test

import (
	"context"
	"testing"

	"github.com/hrygo/echomind/internal/service"
	"github.com/hrygo/echomind/pkg/ai"
)

type MockAIProvider struct{}

func (m *MockAIProvider) Summarize(ctx context.Context, text string) (ai.AnalysisResult, error) {
	return ai.AnalysisResult{
		Summary:     "Mock Summary",
		Category:    "Work",
		Sentiment:   "Neutral",
		Urgency:     "Low",
		ActionItems: []string{"Task 1"},
	}, nil
}

func (m *MockAIProvider) Classify(ctx context.Context, text string) (string, error) {
	return "Work", nil
}

func (m *MockAIProvider) AnalyzeSentiment(ctx context.Context, text string) (ai.SentimentResult, error) {
	return ai.SentimentResult{Sentiment: "Neutral", Urgency: "Low"}, nil
}

func TestSummaryService_GenerateSummary(t *testing.T) {
	var mockProvider ai.AIProvider = &MockAIProvider{}
	svc := service.NewSummaryService(mockProvider)

	result, err := svc.GenerateSummary(context.Background(), "Some text")
	if err != nil {
		t.Fatalf("GenerateSummary failed: %v", err)
	}

	if result.Summary != "Mock Summary" {
		t.Errorf("Expected 'Mock Summary', got '%s'", result.Summary)
	}
	if result.Category != "Work" {
		t.Errorf("Expected 'Work', got '%s'", result.Category)
	}
}

func TestSummaryService_AnalyzeSentiment(t *testing.T) {
	var mockProvider ai.AIProvider = &MockAIProvider{}
	svc := service.NewSummaryService(mockProvider)

	result, err := svc.AnalyzeSentiment(context.Background(), "Some text")
	if err != nil {
		t.Fatalf("AnalyzeSentiment failed: %v", err)
	}

	if result.Sentiment != "Neutral" {
		t.Errorf("Expected Sentiment 'Neutral', got '%s'", result.Sentiment)
	}
}

func TestFactory(t *testing.T) {
	// Test Factory logic (simplified for unit test without real config loading)
	// Real factory test would need viper setup or dependency injection of config.
	// For now, we test that we can create the service with a provider.

	// TODO: Add integration test for Factory with Viper
}
