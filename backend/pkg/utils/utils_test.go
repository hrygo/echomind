package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStripHTML(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple HTML",
			input:    "<p>Hello <b>World</b></p>",
			expected: "Hello World",
		},
		{
			name:     "With Attributes",
			input:    "<div class='container'>Content</div>",
			expected: "Content",
		},
		{
			name:     "With Entities",
			input:    "Fish &amp; Chips",
			expected: "Fish & Chips",
		},
		{
			name:     "Empty",
			input:    "",
			expected: "",
		},
		{
			name:     "Nested Tags",
			input:    "<div><p>Nested</p></div>",
			expected: "Nested",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StripHTML(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestChunker(t *testing.T) {
	chunker := NewTextChunker(10) // ~40 chars

	t.Run("Short Text", func(t *testing.T) {
		text := "Short text."
		chunks := chunker.Chunk(text)
		assert.Len(t, chunks, 1)
		assert.Equal(t, "Short text.", chunks[0])
	})

	t.Run("Multiple Paragraphs", func(t *testing.T) {
		text := "Para 1.\n\nPara 2."
		chunks := chunker.Chunk(text)
		assert.Len(t, chunks, 1) // Should fit in one chunk (7+7 < 40)
		assert.Equal(t, "Para 1.\n\nPara 2.", chunks[0])
	})

	t.Run("Long Paragraph Split", func(t *testing.T) {
		// 50 chars > 40 chars limit
		text := "This is a very long paragraph that should be split into multiple chunks because it exceeds the limit."
		chunks := chunker.Chunk(text)
		assert.Greater(t, len(chunks), 1)
	})
}
