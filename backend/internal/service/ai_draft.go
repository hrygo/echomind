package service

import (
	"context"

	"github.com/hrygo/echomind/pkg/ai"
)

type AIDraftService struct {
	aiProvider ai.AIProvider
}

func NewAIDraftService(aiProvider ai.AIProvider) *AIDraftService {
	return &AIDraftService{aiProvider: aiProvider}
}

func (s *AIDraftService) GenerateDraftReply(ctx context.Context, emailContent, userPrompt string) (string, error) {
	return s.aiProvider.GenerateDraftReply(ctx, emailContent, userPrompt)
}
