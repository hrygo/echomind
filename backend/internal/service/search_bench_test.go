package service_test

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hrygo/echomind/internal/model"
	"github.com/hrygo/echomind/internal/service"
	"github.com/pgvector/pgvector-go"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// MockEmbeddingProvider for benchmarking to avoid API calls
type MockEmbeddingProvider struct{}

func (m *MockEmbeddingProvider) Embed(ctx context.Context, text string) ([]float32, error) {
	// Return a random normalized vector of dimension 1024 (matching current config)
	return randomVector(1024), nil
}

func (m *MockEmbeddingProvider) EmbedBatch(ctx context.Context, texts []string) ([][]float32, error) {
	var results [][]float32
	for range texts {
		results = append(results, randomVector(1024))
	}
	return results, nil
}

func (m *MockEmbeddingProvider) GetDimensions() int {
	return 1024
}

func randomVector(dim int) []float32 {
	vec := make([]float32, dim)
	var sumSq float32
	for i := 0; i < dim; i++ {
		val := rand.Float32()
		vec[i] = val
		sumSq += val * val
	}
	// Normalize (optional, but good for cosine distance)
	norm := float32(1.0) // Simplification: real normalization requires sqrt(sumSq)
	_ = norm             // avoid unused
	return vec
}

func getBenchmarkDB(b *testing.B) *gorm.DB {
	dsn := os.Getenv("TEST_DB_DSN")
	if dsn == "" {
		// Defaults matching docker-compose.yml
		dsn = "host=localhost user=user password=password dbname=echomind_db port=5432 sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Discard, // Silence logs during benchmark
	})
	if err != nil {
		b.Skipf("Skipping benchmark: failed to connect to DB: %v", err)
	}
	return db
}

func setupBenchmarkData(b *testing.B, db *gorm.DB, count int) uuid.UUID {
	userID := uuid.New()

	// Check if we already have enough data for a benchmark user (optimization)
	// For now, we insert fresh data to be accurate about volume, or use a dedicated benchmark user.
	// Inserting 10k rows might take too long for a simple "go test -bench".
	// We will insert in batches.

	b.Logf("Seeding %d emails for benchmark user %s...", count, userID)

	emails := make([]model.Email, 0, 1000)
	embeddings := make([]model.EmailEmbedding, 0, 1000)

	batchSize := 1000
	for i := 0; i < count; i++ {
		emailID := uuid.New()
		emails = append(emails, model.Email{
			ID:        emailID,
			UserID:    userID,
			MessageID: fmt.Sprintf("bench-%s-%d", userID, i),
			Subject:   fmt.Sprintf("Benchmark Email %d", i),
			Snippet:   "This is a benchmark email snippet for search testing.",
			Date:      time.Now(),
			BodyText:  "Full body text would go here...",
		})

		embeddings = append(embeddings, model.EmailEmbedding{
			EmailID: emailID,
			Vector:  pgvector.NewVector(randomVector(1536)),
		})

		if len(emails) >= batchSize {
			if err := db.Create(&emails).Error; err != nil {
				b.Fatalf("Failed to create emails: %v", err)
			}
			if err := db.Create(&embeddings).Error; err != nil {
				b.Fatalf("Failed to create embeddings: %v", err)
			}
			emails = emails[:0]
			embeddings = embeddings[:0]
		}
	}
	// Insert remaining
	if len(emails) > 0 {
		if err := db.Create(&emails).Error; err != nil {
			b.Fatalf("Failed to create emails: %v", err)
		}
		if err := db.Create(&embeddings).Error; err != nil {
			b.Fatalf("Failed to create embeddings: %v", err)
		}
	}

	return userID
}

func cleanupBenchmarkData(db *gorm.DB, userID uuid.UUID) {
	// Delete embeddings first (cascade might handle it, but being explicit)
	// Accessing table via model
	// We need to find emails for this user
	db.Exec("DELETE FROM email_embeddings WHERE email_id IN (SELECT id FROM emails WHERE user_id = ?)", userID)
	db.Exec("DELETE FROM emails WHERE user_id = ?", userID)
}

func BenchmarkSearch_100(b *testing.B)  { benchmarkSearch(b, 100) }
func BenchmarkSearch_1000(b *testing.B) { benchmarkSearch(b, 1000) }

// func BenchmarkSearch_10000(b *testing.B) { benchmarkSearch(b, 10000) } // Uncomment for heavier tests

func benchmarkSearch(b *testing.B, count int) {
	db := getBenchmarkDB(b)

	// Ensure vector extension
	db.Exec("CREATE EXTENSION IF NOT EXISTS vector")

	// Ensure schema
	if err := db.AutoMigrate(&model.Email{}, &model.EmailEmbedding{}); err != nil {
		b.Fatalf("Failed to migrate: %v", err)
	}

	// Create HNSW index manually to ensure correct operator class
	// Note: This might fail if data exists and is not normalized, but with randomVector it's fine.
	// We drop it first to ensure we test the index creation/usage, or just IF NOT EXISTS.
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_email_embeddings_vector ON email_embeddings USING hnsw (vector vector_cosine_ops)").Error; err != nil {
		b.Logf("Warning: Failed to create HNSW index: %v. Proceeding without it (performance might suffer).", err)
	}

	userID := setupBenchmarkData(b, db, count)
	defer cleanupBenchmarkData(db, userID)

	svc := service.NewSearchService(db, &MockEmbeddingProvider{}, nil)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := svc.Search(ctx, userID, "benchmark query", service.SearchFilters{}, 10)
		if err != nil {
			b.Fatalf("Search failed: %v", err)
		}
	}
}
