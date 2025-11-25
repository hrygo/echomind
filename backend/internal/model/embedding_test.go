package model

import (
	"testing"

	"github.com/google/uuid"
	"github.com/pgvector/pgvector-go"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestEmailEmbeddingVectorConversion(t *testing.T) {
	// Use in-memory SQLite for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	// Migrate the schema
	err = db.AutoMigrate(&EmailEmbedding{})
	if err != nil {
		t.Fatalf("Failed to migrate schema: %v", err)
	}

	tests := []struct {
		name         string
		inputVector  []float32
		expectedDims int
	}{
		{
			name:         "768 dimensions (Gemini)",
			inputVector:  make([]float32, 768),
			expectedDims: 768,
		},
		{
			name:         "1024 dimensions (SiliconFlow)",
			inputVector:  make([]float32, 1024),
			expectedDims: 1024,
		},
		{
			name:         "1536 dimensions (OpenAI)",
			inputVector:  make([]float32, 1536),
			expectedDims: 1536,
		},
		{
			name:         "Short vector needs padding",
			inputVector:  make([]float32, 512),
			expectedDims: 512,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set some values in the vector for testing
			for i := range tt.inputVector {
				tt.inputVector[i] = float32(i + 1)
			}

			embedding := EmailEmbedding{
				EmailID:   uuid.New(),
				Content:   "test content",
				Vector:    pgvector.NewVector(tt.inputVector),
				Dimensions: len(tt.inputVector),
			}

			// BeforeCreate should be called automatically
			err := db.Create(&embedding).Error
			if err != nil {
				t.Fatalf("Failed to create embedding: %v", err)
			}

			// Verify the vector was converted to 1536 dimensions
			vectorSlice := embedding.Vector.Slice()
			if len(vectorSlice) != 1536 {
				t.Errorf("Expected vector length 1536, got %d", len(vectorSlice))
			}

			// Verify the Dimensions field tracks original size
			if embedding.Dimensions != tt.expectedDims {
				t.Errorf("Expected dimensions %d, got %d", tt.expectedDims, embedding.Dimensions)
			}

			// Verify padding is zero for padded vectors
			if len(tt.inputVector) < 1536 {
				for i := len(tt.inputVector); i < 1536; i++ {
					if vectorSlice[i] != 0 {
						t.Errorf("Expected zero padding at index %d, got %f", i, vectorSlice[i])
					}
				}
			}

			// Verify original values are preserved
			for i, expected := range tt.inputVector {
				if vectorSlice[i] != expected {
					t.Errorf("Expected %f at index %d, got %f", expected, i, vectorSlice[i])
				}
			}
		})
	}
}

func TestEmailEmbeddingTableName(t *testing.T) {
	embedding := EmailEmbedding{}
	expected := "email_embeddings"
	if embedding.TableName() != expected {
		t.Errorf("Expected table name %s, got %s", expected, embedding.TableName())
	}
}