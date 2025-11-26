package model

import (
	"testing"

	"github.com/pgvector/pgvector-go"
)

func TestEmailEmbeddingTableName(t *testing.T) {
	embedding := EmailEmbedding{}
	expected := "email_embeddings"
	if embedding.TableName() != expected {
		t.Errorf("Expected table name %s, got %s", expected, embedding.TableName())
	}
}

func TestEmailEmbeddingFixedDimensions(t *testing.T) {
	// Test that embeddings maintain fixed 1024 dimensions (规约化 approach)
	testVector := make([]float32, 1024)
	for i := range testVector {
		testVector[i] = float32(i + 1)
	}

	embedding := EmailEmbedding{
		Vector:     pgvector.NewVector(testVector),
		Dimensions: 1024, // Always 1024 with规约化 approach
	}

	// Verify dimensions are fixed
	if embedding.Dimensions != 1024 {
		t.Errorf("Expected fixed dimensions 1024, got %d", embedding.Dimensions)
	}

	// Verify vector length matches
	vectorSlice := embedding.Vector.Slice()
	if len(vectorSlice) != 1024 {
		t.Errorf("Expected vector length 1024, got %d", len(vectorSlice))
	}
}
