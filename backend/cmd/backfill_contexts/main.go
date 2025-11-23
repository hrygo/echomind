package main

import (
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/hrygo/echomind/internal/model"
	"github.com/hrygo/echomind/internal/service"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Initialize Viper
	vip := viper.New()
	vip.SetConfigFile("configs/config.yaml")
	vip.AddConfigPath(".")
	vip.AutomaticEnv()
	vip.SetEnvPrefix("ECHOMIND")
	vip.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := vip.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	// Database
	dsn := vip.GetString("database.dsn")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	contextService := service.NewContextService(db)

	// Fetch all emails
	var emails []model.Email
	if err := db.Find(&emails).Error; err != nil {
		log.Fatalf("Failed to fetch emails: %v", err)
	}

	log.Printf("Found %d emails to scan for contexts", len(emails))

	success := 0
	matched := 0
	failed := 0

	for _, email := range emails {
		matches, err := contextService.MatchContexts(&email)
		if err != nil {
			log.Printf("Failed to match context for email %s: %v", email.ID, err)
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
				log.Printf("Failed to assign contexts to email %s: %v", email.ID, err)
				failed++
			} else {
				log.Printf("Matched email %s to contexts: %v", email.ID, names)
				matched++
			}
		}
		success++
	}

	log.Printf("Backfill complete. Scanned: %d, Matched: %d, Failed: %d", success, matched, failed)
}
