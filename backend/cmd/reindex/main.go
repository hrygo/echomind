package main

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/hibiken/asynq"
	"github.com/hrygo/echomind/internal/model"
	"github.com/hrygo/echomind/internal/tasks"
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

	// Asynq Client
	redisAddr := vip.GetString("redis.addr")
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
	defer client.Close()

	// Fetch all emails
	// In a real production scenario, we would paginate this.
	var emails []model.Email
	if err := db.Find(&emails).Error; err != nil {
		log.Fatalf("Failed to fetch emails: %v", err)
	}

	log.Printf("Found %d emails to reindex", len(emails))

	for _, email := range emails {
		// Enqueue analysis task
		// This will trigger summary + embedding generation
		// Note: This might re-summarize emails, which costs tokens.
		// If we only want to generate embeddings, we should create a separate task or flag.
		// For now, per plan, we reuse the analyze task.

		payload, err := json.Marshal(tasks.EmailAnalyzePayload{
			EmailID: email.ID,
			UserID:  email.UserID,
		})
		if err != nil {
			log.Printf("Failed to marshal payload for email %s: %v", email.ID, err)
			continue
		}

		task := asynq.NewTask(tasks.TypeEmailAnalyze, payload)
		info, err := client.Enqueue(task, asynq.MaxRetry(3), asynq.Timeout(5*time.Minute))
		if err != nil {
			log.Printf("Failed to enqueue task for email %s: %v", email.ID, err)
			continue
		}

		log.Printf("Enqueued task %s for email %s", info.ID, email.ID)
	}

	log.Println("Reindex complete (tasks enqueued)")
}
