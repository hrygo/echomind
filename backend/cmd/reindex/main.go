package main

import (
	"context"
	"log"

	"github.com/hrygo/echomind/internal/app"
	"github.com/hrygo/echomind/internal/model"
	"github.com/hrygo/echomind/pkg/logger"
)

func main() {
	// Parse CLI configuration
	cli := app.ParseCLI()

	// Initialize application container
	container, err := app.NewContainer(cli.ConfigPath, cli.IsProduction)
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}
	defer container.Close()

	// Fetch all emails
	var emails []model.Email
	if err := container.DB.Select("id, subject, snippet, body_text").Find(&emails).Error; err != nil {
		container.Logger.Fatal("Failed to fetch emails", logger.Error(err))
	}

	container.Logger.Info("Starting reindex process",
		logger.Int("email_count", len(emails)))

	ctx := context.Background()
	success := 0
	failed := 0

	for _, email := range emails {
		container.Logger.Info("Reindexing email",
			logger.String("email_id", email.ID.String()),
			logger.String("subject", email.Subject))
		if err := container.SearchService.GenerateAndSaveEmbedding(ctx, &email, container.ChunkSize()); err != nil {
			container.Logger.Warn("Failed to reindex email",
				logger.String("email_id", email.ID.String()),
				logger.Error(err))
			failed++
		} else {
			success++
		}
	}

	container.Logger.Info("Reindex complete",
		logger.Int("success", success),
		logger.Int("failed", failed))
}
