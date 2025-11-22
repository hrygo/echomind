package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hrygo/echomind/pkg/ai"
	"github.com/pgvector/pgvector-go"
	"gorm.io/gorm"
)

type SearchService struct {
	db       *gorm.DB
	embedder ai.EmbeddingProvider
}

func NewSearchService(db *gorm.DB, embedder ai.EmbeddingProvider) *SearchService {
	return &SearchService{
		db:       db,
		embedder: embedder,
	}
}

type SearchResult struct {
	EmailID uuid.UUID `json:"email_id"`
	Subject string    `json:"subject"`
	Snippet string    `json:"snippet"`
	Sender  string    `json:"sender"`
	Date    time.Time `json:"date"`
	Score   float64   `json:"score"` // Similarity score (1 - distance)
}

type SearchFilters struct {
	Sender    string
	StartDate *time.Time
	EndDate   *time.Time
}

func (s *SearchService) Search(ctx context.Context, userID uuid.UUID, query string, filters SearchFilters, limit int) ([]SearchResult, error) {
	// 1. Generate embedding for the query
	queryVector, err := s.embedder.Embed(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to embed query: %w", err)
	}

	// 2. Perform vector search using raw SQL
	// Note: <=> is cosine distance operator in pgvector
	// We join with emails table to filter by user_id and get metadata
	var results []SearchResult

	// Base query
	sql := `
		SELECT 
			e.id as email_id,
			e.subject,
			e.snippet,
			e.sender,
			e.date,
			1 - (ee.vector <=> ?) as score
		FROM email_embeddings ee
		JOIN emails e ON e.id = ee.email_id
		WHERE e.user_id = ?
	`
	args := []interface{}{pgvector.NewVector(queryVector), userID}

	// Apply filters
	if filters.Sender != "" {
		sql += " AND e.sender ILIKE ?"
		args = append(args, "%"+filters.Sender+"%")
	}
	if filters.StartDate != nil {
		sql += " AND e.date >= ?"
		args = append(args, filters.StartDate)
	}
	if filters.EndDate != nil {
		sql += " AND e.date <= ?"
		args = append(args, filters.EndDate)
	}

	// Order by and Limit
	sql += " ORDER BY ee.vector <=> ? LIMIT ?"
	args = append(args, pgvector.NewVector(queryVector), limit)

	err = s.db.WithContext(ctx).Raw(sql, args...).Scan(&results).Error

	if err != nil {
		return nil, fmt.Errorf("search query failed: %w", err)
	}

	return results, nil
}
