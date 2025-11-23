package main

import (
	"context"
	"log"

	"github.com/hrygo/echomind/internal/bootstrap"
	"github.com/hrygo/echomind/internal/model"
	"github.com/hrygo/echomind/internal/service"
	"github.com/hrygo/echomind/pkg/ai"
)

func main() {
	// 1. Bootstrap
	app, err := bootstrap.Init("configs/config.yaml", false)
	if err != nil {
		log.Fatalf("Bootstrap failed: %v", err)
	}
	defer app.Close()

	// 2. Services
	aiProvider, err := service.NewAIProvider(&app.Config.AI)
	if err != nil {
		app.Sugar.Fatalf("Failed to create AI provider: %v", err)
	}

	embedder, ok := aiProvider.(ai.EmbeddingProvider)
	if !ok {
		app.Sugar.Fatal("AI provider does not implement EmbeddingProvider")
	}

	searchService := service.NewSearchService(app.DB, embedder)
	chunkSize := 1000 // Default, should be in config

	// 3. Logic
	var emails []model.Email
	if err := app.DB.Select("id, subject, snippet, body_text").Find(&emails).Error; err != nil {
		app.Sugar.Fatalf("Failed to fetch emails: %v", err)
	}

	app.Sugar.Infof("Found %d emails to reindex", len(emails))

	ctx := context.Background()
	success := 0
	failed := 0

	for _, email := range emails {
		app.Sugar.Infof("Reindexing email %s: %s", email.ID, email.Subject)
		if err := searchService.GenerateAndSaveEmbedding(ctx, &email, chunkSize); err != nil {
			app.Sugar.Warnf("Failed to reindex email %s: %v", email.ID, err)
			failed++
		} else {
			success++
		}
	}

	app.Sugar.Infof("Reindex complete. Success: %d, Failed: %d", success, failed)
}
