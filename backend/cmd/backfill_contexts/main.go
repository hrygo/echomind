package main

import (
	"log"

	"github.com/google/uuid"
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
	if err := container.DB.Find(&emails).Error; err != nil {
		container.Logger.Fatal("Failed to fetch emails", logger.Error(err))
	}

	container.Logger.Info("Starting context backfill process",
		logger.Int("email_count", len(emails)))

	success := 0
	matched := 0
	failed := 0

	for _, email := range emails {
		matches, err := container.ContextService.MatchContexts(&email)
		if err != nil {
			container.Logger.Warn("Failed to match context for email",
				logger.String("email_id", email.ID.String()),
				logger.Error(err))
			failed++
			continue
		}

		if len(matches) > 0 {
			var contextIDs []uuid.UUID
			names := []string{}
			for _, m := range matches {
				contextIDs = append(contextIDs, m.ID)
				names = append(names, m.Name)
			}

			if err := container.ContextService.AssignContextsToEmail(email.ID, contextIDs); err != nil {
				container.Logger.Warn("Failed to assign contexts to email",
					logger.String("email_id", email.ID.String()),
					logger.Error(err))
				failed++
			} else {
				container.Logger.Info("Matched email to contexts",
					logger.String("email_id", email.ID.String()),
					logger.Strings("context_names", names))
				matched++
			}
		}
		success++
	}

	container.Logger.Info("Context backfill complete",
		logger.Int("scanned", success),
		logger.Int("matched", matched),
		logger.Int("failed", failed))
}
