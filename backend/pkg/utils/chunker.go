package utils

import (
	"strings"
)

// TextChunker handles splitting text into smaller chunks suitable for embedding.
type TextChunker struct {
	MaxTokens int
}

// NewTextChunker creates a new chunker with the specified max tokens per chunk.
// A rough approximation is 1 token ~= 4 chars.
func NewTextChunker(maxTokens int) *TextChunker {
	return &TextChunker{
		MaxTokens: maxTokens,
	}
}

// Chunk splits the text into chunks.
// It uses a simple paragraph-based strategy first, then splits by sentences if needed.
func (c *TextChunker) Chunk(text string) []string {
	if text == "" {
		return []string{}
	}

	// Normalize newlines
	text = strings.ReplaceAll(text, "\r\n", "\n")

	// Split by paragraphs
	paragraphs := strings.Split(text, "\n\n")

	var chunks []string
	var currentChunk strings.Builder

	// Rough char limit based on tokens (4 chars per token)
	// We use a safety margin of 10%
	maxChars := int(float64(c.MaxTokens) * 4 * 0.9)

	for _, p := range paragraphs {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}

		// If a single paragraph is too long, we need to split it further (e.g. by sentences)
		// For now, we'll just truncate/split hard if it exceeds the limit,
		// but a better approach would be sentence splitting.
		if len(p) > maxChars {
			// If current chunk is not empty, flush it
			if currentChunk.Len() > 0 {
				chunks = append(chunks, currentChunk.String())
				currentChunk.Reset()
			}

			// Split long paragraph
			subChunks := c.splitLongText(p, maxChars)
			chunks = append(chunks, subChunks...)
			continue
		}

		// Check if adding this paragraph exceeds the limit
		if currentChunk.Len()+len(p)+2 > maxChars {
			chunks = append(chunks, currentChunk.String())
			currentChunk.Reset()
		}

		if currentChunk.Len() > 0 {
			currentChunk.WriteString("\n\n")
		}
		currentChunk.WriteString(p)
	}

	if currentChunk.Len() > 0 {
		chunks = append(chunks, currentChunk.String())
	}

	return chunks
}

func (c *TextChunker) splitLongText(text string, limit int) []string {
	var chunks []string
	for len(text) > limit {
		// Try to find a split point (newline or space) near the limit
		splitIdx := strings.LastIndexAny(text[:limit], ".!?\n ")
		if splitIdx == -1 {
			splitIdx = limit // Hard split if no natural break found
		} else {
			splitIdx++ // Include the punctuation/space
		}

		chunks = append(chunks, strings.TrimSpace(text[:splitIdx]))
		text = text[splitIdx:]
	}
	if len(text) > 0 {
		chunks = append(chunks, strings.TrimSpace(text))
	}
	return chunks
}
