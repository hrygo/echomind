package service

import (
	"context"

	"github.com/hrygo/echomind/pkg/ai"
)

type SummaryService struct {
	provider ai.AIProvider
}

func NewSummaryService(provider ai.AIProvider) *SummaryService {
	return &SummaryService{
		provider: provider,
	}
}

func (s *SummaryService) GenerateSummary(ctx context.Context, text string) (ai.AnalysisResult, error) {
	return s.provider.Summarize(ctx, text)
}

func (s *SummaryService) AnalyzeSentiment(ctx context.Context, text string) (ai.SentimentResult, error) {
	return s.provider.AnalyzeSentiment(ctx, text)
}
