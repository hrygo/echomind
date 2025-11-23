package main

import (
	"context"
	"log"

	"github.com/hrygo/echomind/internal/app"
	"github.com/hrygo/echomind/internal/model"
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
		container.Sugar.Fatalf("Failed to fetch emails: %v", err)
	}

	container.Sugar.Infof("Found %d emails to reindex", len(emails))

	ctx := context.Background()
	success := 0
	failed := 0

	for _, email := range emails {
		container.Sugar.Infof("Reindexing email %s: %s", email.ID, email.Subject)
		if err := container.SearchService.GenerateAndSaveEmbedding(ctx, &email, container.ChunkSize()); err != nil {
			container.Sugar.Warnf("Failed to reindex email %s: %v", email.ID, err)
			failed++
		} else {
			success++
		}
	}

	container.Sugar.Infof("Reindex complete. Success: %d, Failed: %d", success, failed)
}
