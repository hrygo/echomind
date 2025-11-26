package service_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hrygo/echomind/internal/model"
	"github.com/hrygo/echomind/internal/service"
	"github.com/pgvector/pgvector-go"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Use MockEmbeddingProvider from search_bench_test.go (assumed shared package)
// If not shared, we redefine or move it to a common test file.
// Since they are in the same package `service_test`, they are shared.

func getIntegrationDB(t *testing.T) *gorm.DB {
	dsn := os.Getenv("TEST_DB_DSN")
	if dsn == "" {
		// Default to local docker instance
		dsn = "host=localhost user=user password=password dbname=echomind_db port=5432 sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Logf("Skipping integration test: failed to connect to DB: %v", err)
		return nil
	}

	// Check connectivity
	sqlDB, err := db.DB()
	if err != nil || sqlDB.Ping() != nil {
		t.Logf("Skipping integration test: DB unreachable")
		return nil
	}

	return db
}

func TestSearchService_Search_Integration(t *testing.T) {
	db := getIntegrationDB(t)
	if db == nil {
		t.Skip("Skipping integration test (requires Postgres with pgvector)")
	}

	// Cleanup & Setup
	// Ensure vector extension
	db.Exec("CREATE EXTENSION IF NOT EXISTS vector")

	// Ensure schema
	err := db.AutoMigrate(&model.Email{}, &model.EmailEmbedding{})
	assert.NoError(t, err)

	userID := uuid.New()

	// Cleanup
	defer func() {
		db.Exec("DELETE FROM email_embeddings WHERE email_id IN (SELECT id FROM emails WHERE user_id = ?)", userID)
		db.Exec("DELETE FROM emails WHERE user_id = ?", userID)
	}()

	// Seed Data
	emailID := uuid.New()
	email := model.Email{
		ID:        emailID,
		UserID:    userID,
		MessageID: "test-msg-1",
		Subject:   "Project Alpha Kickoff",
		Snippet:   "Meeting to discuss Q4 goals.",
		Date:      time.Now(),
		BodyText:  "Full body text...",
	}
	err = db.Create(&email).Error
	assert.NoError(t, err)

	// Create Embedding (Mocked vector)
	// Query: "Project" -> Vector A
	// Email: "Project Alpha" -> Vector B (Close to A)
	// We simulate this by inserting a vector and querying with a similar one.

	// Create specific vector
	embeddingVec := make([]float32, 1024)
	embeddingVec[0] = 1.0

	embedding := model.EmailEmbedding{
		EmailID: emailID,
		Vector:  pgvector.NewVector(embeddingVec),
	}
	err = db.Create(&embedding).Error
	assert.NoError(t, err)

	// Mock Embedder that returns the SAME vector for the query
	// This ensures distance is 0 (Score 1.0)
	fixedEmbedder := &MockFixedEmbedder{Vector: embeddingVec}

	svc := service.NewSearchService(db, fixedEmbedder, nil)

	// Test Search
	results, err := svc.Search(context.Background(), userID, "Project", service.SearchFilters{}, 10)
	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, emailID, results[0].EmailID)
	assert.InDelta(t, 1.0, results[0].Score, 0.0001) // Score should be 1 - distance (0) = 1
}

type MockFixedEmbedder struct {
	Vector []float32
}

func (m *MockFixedEmbedder) Embed(ctx context.Context, text string) ([]float32, error) {
	return m.Vector, nil
}

func (m *MockFixedEmbedder) EmbedBatch(ctx context.Context, texts []string) ([][]float32, error) {
	var results [][]float32
	for range texts {
		results = append(results, m.Vector)
	}
	return results, nil
}

func (m *MockFixedEmbedder) GetDimensions() int {
	return len(m.Vector)
}
