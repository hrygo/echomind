package main

import (
	"log"

	"github.com/google/uuid"
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
	if err := container.DB.Find(&emails).Error; err != nil {
		container.Sugar.Fatalf("Failed to fetch emails: %v", err)
	}

	container.Sugar.Infof("Found %d emails to scan for contexts", len(emails))

	success := 0
	matched := 0
	failed := 0

	for _, email := range emails {
		matches, err := container.ContextService.MatchContexts(&email)
		if err != nil {
			container.Sugar.Warnf("Failed to match context for email %s: %v", email.ID, err)
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
				container.Sugar.Warnf("Failed to assign contexts to email %s: %v", email.ID, err)
				failed++
			} else {
				container.Sugar.Infof("Matched email %s to contexts: %v", email.ID, names)
				matched++
			}
		}
		success++
	}

	container.Sugar.Infof("Backfill complete. Scanned: %d, Matched: %d, Failed: %d", success, matched, failed)
}
