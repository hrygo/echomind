package main

import (
	"log"

	"github.com/google/uuid"
	"github.com/hrygo/echomind/internal/bootstrap"
	"github.com/hrygo/echomind/internal/model"
	"github.com/hrygo/echomind/internal/service"
)

func main() {
	// 1. Bootstrap
	app, err := bootstrap.Init("configs/config.yaml", false)
	if err != nil {
		log.Fatalf("Bootstrap failed: %v", err)
	}
	defer app.Close()

	contextService := service.NewContextService(app.DB)

	// Fetch all emails
	var emails []model.Email
	if err := app.DB.Find(&emails).Error; err != nil {
		app.Sugar.Fatalf("Failed to fetch emails: %v", err)
	}

	app.Sugar.Infof("Found %d emails to scan for contexts", len(emails))

	success := 0
	matched := 0
	failed := 0

	for _, email := range emails {
		matches, err := contextService.MatchContexts(&email)
		if err != nil {
			app.Sugar.Warnf("Failed to match context for email %s: %v", email.ID, err)
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
			
			if err := contextService.AssignContextsToEmail(email.ID, contextIDs); err != nil {
				app.Sugar.Warnf("Failed to assign contexts to email %s: %v", email.ID, err)
				failed++
			} else {
				app.Sugar.Infof("Matched email %s to contexts: %v", email.ID, names)
				matched++
			}
		}
		success++
	}

	app.Sugar.Infof("Backfill complete. Scanned: %d, Matched: %d, Failed: %d", success, matched, failed)
}