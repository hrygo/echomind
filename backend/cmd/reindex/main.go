package main

import (
	"context"
	"log"
	"strings"

	"github.com/hrygo/echomind/configs"
	"github.com/hrygo/echomind/internal/model"
	"github.com/hrygo/echomind/internal/service"
	"github.com/hrygo/echomind/pkg/ai"
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

	// Load entire config into struct
	var appConfig configs.Config
	if err := vip.Unmarshal(&appConfig); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	// Database
	dsn := vip.GetString("database.dsn")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// AI Service
	aiProvider, err := service.NewAIProvider(&appConfig.AI)
	if err != nil {
		log.Fatalf("Failed to create AI provider: %v", err)
	}

	embedder, ok := aiProvider.(ai.EmbeddingProvider)
	if !ok {
		log.Fatal("AI provider does not implement EmbeddingProvider")
	}

	searchService := service.NewSearchService(db, embedder)
	chunkSize := vip.GetInt("ai.chunk_size")
	if chunkSize <= 0 {
		chunkSize = 1000
	}

	// Fetch all emails
	// For large datasets, use batching/pagination.
	var emails []model.Email
	// We can filter by user_id or reindex all.
	// Let's reindex all for now.
	// Optimize: only fetch ID, Subject, Snippet, BodyText
	if err := db.Select("id, subject, snippet, body_text").Find(&emails).Error; err != nil {
		log.Fatalf("Failed to fetch emails: %v", err)
	}

	log.Printf("Found %d emails to reindex", len(emails))

	ctx := context.Background()
	success := 0
	failed := 0

	for _, email := range emails {
		log.Printf("Reindexing email %s: %s", email.ID, email.Subject)
		if err := searchService.GenerateAndSaveEmbedding(ctx, &email, chunkSize); err != nil {
			log.Printf("Failed to reindex email %s: %v", email.ID, err)
			failed++
		} else {
			success++
		}
	}

	log.Printf("Reindex complete. Success: %d, Failed: %d", success, failed)
}