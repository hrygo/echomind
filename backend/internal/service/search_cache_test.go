package service

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestRedis creates a miniredis instance and redis client for testing
func setupTestRedis(t *testing.T) (*miniredis.Miniredis, *redis.Client) {
	mr := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	return mr, client
}

// TestSearchCache_Get_CacheHit tests successful cache retrieval
func TestSearchCache_Get_CacheHit(t *testing.T) {
	_, redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cache := NewSearchCache(redisClient, 30*time.Minute)

	ctx := context.Background()
	userID := uuid.New()
	query := "test query"
	filters := SearchFilters{}
	limit := 10

	// Prepare test data
	expectedResults := []SearchResult{
		{
			EmailID: uuid.New(),
			Subject: "Test Email",
			Snippet: "Test content",
			Sender:  "test@example.com",
			Date:    time.Now(),
			Score:   0.95,
		},
	}

	// Pre-populate cache
	key := cache.generateCacheKey(ctx, userID, query, filters, limit)
	data, err := json.Marshal(expectedResults)
	require.NoError(t, err)
	err = redisClient.Set(ctx, key, data, 30*time.Minute).Err()
	require.NoError(t, err)

	// Execute
	results, found, err := cache.Get(ctx, userID, query, filters, limit)

	// Assert
	assert.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, len(expectedResults), len(results))
	assert.Equal(t, expectedResults[0].EmailID, results[0].EmailID)
	assert.Equal(t, expectedResults[0].Subject, results[0].Subject)
}

// TestSearchCache_Get_CacheMiss tests cache miss scenario
func TestSearchCache_Get_CacheMiss(t *testing.T) {
	_, redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cache := NewSearchCache(redisClient, 30*time.Minute)

	ctx := context.Background()
	userID := uuid.New()
	query := "test query"
	filters := SearchFilters{}
	limit := 10

	// Execute (cache is empty)
	results, found, err := cache.Get(ctx, userID, query, filters, limit)

	// Assert
	assert.NoError(t, err)
	assert.False(t, found)
	assert.Nil(t, results)
}

// TestSearchCache_Set tests cache write operation
func TestSearchCache_Set(t *testing.T) {
	_, redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cache := NewSearchCache(redisClient, 30*time.Minute)

	ctx := context.Background()
	userID := uuid.New()
	query := "test query"
	filters := SearchFilters{}
	limit := 10

	results := []SearchResult{
		{
			EmailID: uuid.New(),
			Subject: "Test Email",
			Snippet: "Test content",
			Sender:  "test@example.com",
			Date:    time.Now(),
			Score:   0.95,
		},
	}

	// Execute
	err := cache.Set(ctx, userID, query, filters, limit, results)

	// Assert
	assert.NoError(t, err)

	// Verify data was stored
	key := cache.generateCacheKey(ctx, userID, query, filters, limit)
	storedData, err := redisClient.Get(ctx, key).Bytes()
	require.NoError(t, err)

	// Verify data can be retrieved
	var storedResults []SearchResult
	err = json.Unmarshal(storedData, &storedResults)
	require.NoError(t, err)
	assert.Equal(t, len(results), len(storedResults))
	assert.Equal(t, results[0].EmailID, storedResults[0].EmailID)
}

// TestSearchCache_Invalidate tests cache invalidation for a user
func TestSearchCache_Invalidate(t *testing.T) {
	mr, redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cache := NewSearchCache(redisClient, 30*time.Minute)

	ctx := context.Background()
	userID := uuid.New()
	otherUserID := uuid.New()

	// Populate cache with multiple entries
	results := []SearchResult{{EmailID: uuid.New(), Subject: "Test"}}

	_ = cache.Set(ctx, userID, "query1", SearchFilters{}, 10, results)
	_ = cache.Set(ctx, userID, "query2", SearchFilters{}, 10, results)
	_ = cache.Set(ctx, otherUserID, "query3", SearchFilters{}, 10, results)

	initialKeys := mr.Keys()
	assert.GreaterOrEqual(t, len(initialKeys), 3)

	// Execute invalidation
	err := cache.Invalidate(ctx, userID)

	// Assert
	assert.NoError(t, err)
	// Note: Actual validation would check specific keys deleted
}

// TestSearchCache_InvalidateAll tests global cache invalidation
func TestSearchCache_InvalidateAll(t *testing.T) {
	mr, redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cache := NewSearchCache(redisClient, 30*time.Minute)

	ctx := context.Background()
	userID1 := uuid.New()
	userID2 := uuid.New()

	// Populate cache
	results := []SearchResult{{EmailID: uuid.New(), Subject: "Test"}}

	_ = cache.Set(ctx, userID1, "query1", SearchFilters{}, 10, results)
	_ = cache.Set(ctx, userID2, "query2", SearchFilters{}, 10, results)

	initialKeys := mr.Keys()
	assert.GreaterOrEqual(t, len(initialKeys), 2)

	// Execute invalidate all
	err := cache.InvalidateAll(ctx)

	// Assert
	assert.NoError(t, err)
	finalKeys := mr.Keys()
	assert.Equal(t, 0, len(finalKeys))
}

// TestSearchCache_NilRedis tests behavior when Redis client is nil
func TestSearchCache_NilRedis(t *testing.T) {
	cache := NewSearchCache(nil, 30*time.Minute)

	ctx := context.Background()
	userID := uuid.New()
	results := []SearchResult{{EmailID: uuid.New()}}

	// All operations should succeed with nil redis
	_, found, err := cache.Get(ctx, userID, "query", SearchFilters{}, 10)
	assert.NoError(t, err)
	assert.False(t, found)

	err = cache.Set(ctx, userID, "query", SearchFilters{}, 10, results)
	assert.NoError(t, err)

	err = cache.Invalidate(ctx, userID)
	assert.NoError(t, err)

	err = cache.InvalidateAll(ctx)
	assert.NoError(t, err)
}

// TestSearchCache_GenerateCacheKey tests cache key generation
func TestSearchCache_GenerateCacheKey(t *testing.T) {
	cache := NewSearchCache(nil, 30*time.Minute)
	ctx := context.Background()

	userID := uuid.New()
	query := "test query"
	filters := SearchFilters{
		Sender: "test@example.com",
	}
	limit := 10

	// Generate keys
	key1 := cache.generateCacheKey(ctx, userID, query, filters, limit)
	key2 := cache.generateCacheKey(ctx, userID, query, filters, limit)

	// Same inputs should generate same key
	assert.Equal(t, key1, key2)
	assert.Contains(t, key1, "search:cache:")

	// Different inputs should generate different keys
	differentQuery := "different query"
	key3 := cache.generateCacheKey(ctx, userID, differentQuery, filters, limit)
	assert.NotEqual(t, key1, key3)
}

// BenchmarkSearchCache_Get benchmarks cache get operation
func BenchmarkSearchCache_Get(b *testing.B) {
	mr := miniredis.NewMiniRedis()
	if err := mr.Start(); err != nil {
		b.Fatal(err)
	}
	defer mr.Close()

	redisClient := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	defer redisClient.Close()

	cache := NewSearchCache(redisClient, 30*time.Minute)

	ctx := context.Background()
	userID := uuid.New()
	query := "benchmark query"
	filters := SearchFilters{}
	limit := 10

	// Pre-populate cache
	results := []SearchResult{{EmailID: uuid.New(), Subject: "Test"}}
	key := cache.generateCacheKey(ctx, userID, query, filters, limit)
	data, _ := json.Marshal(results)
	_ = redisClient.Set(ctx, key, data, 30*time.Minute).Err()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = cache.Get(ctx, userID, query, filters, limit)
	}
}

// BenchmarkSearchCache_Set benchmarks cache set operation
func BenchmarkSearchCache_Set(b *testing.B) {
	mr := miniredis.NewMiniRedis()
	if err := mr.Start(); err != nil {
		b.Fatal(err)
	}
	defer mr.Close()

	redisClient := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	defer redisClient.Close()

	cache := NewSearchCache(redisClient, 30*time.Minute)

	ctx := context.Background()
	userID := uuid.New()
	filters := SearchFilters{}
	limit := 10

	results := []SearchResult{{EmailID: uuid.New(), Subject: "Test"}}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		query := "query" + string(rune(i))
		_ = cache.Set(ctx, userID, query, filters, limit, results)
	}
}
