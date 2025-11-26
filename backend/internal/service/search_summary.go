package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/hrygo/echomind/pkg/ai"
)

// SearchSummaryService generates AI-powered summaries for search results
type SearchSummaryService struct {
	aiProvider ai.AIProvider
}

// NewSearchSummaryService creates a new search summary service
func NewSearchSummaryService(aiProvider ai.AIProvider) *SearchSummaryService {
	return &SearchSummaryService{
		aiProvider: aiProvider,
	}
}

// SearchResultsSummary contains the AI-generated summary
type SearchResultsSummary struct {
	NaturalSummary  string   `json:"natural_summary"`
	KeyTopics       []string `json:"key_topics"`
	ImportantPeople []string `json:"important_people"`
	UrgentCount     int      `json:"urgent_count"`
	ActionItems     []string `json:"action_items,omitempty"`
}

// GenerateSummary creates an AI-powered summary of search results
func (s *SearchSummaryService) GenerateSummary(ctx context.Context, results []SearchResult, query string) (*SearchResultsSummary, error) {
	if len(results) == 0 {
		return &SearchResultsSummary{
			NaturalSummary: "未找到相关邮件。",
		}, nil
	}

	// Prepare context for AI
	prompt := s.buildSummaryPrompt(results, query)

	// Call AI provider using Summarize method
	// The Summarize method is designed for text analysis and summary generation
	analysisResult, err := s.aiProvider.Summarize(ctx, prompt)

	if err != nil {
		return nil, fmt.Errorf("AI summary generation failed: %w", err)
	}

	// Parse AI response and extract structured data
	// Use the summary from analysis result
	summary := s.parseSummaryResponse(analysisResult.Summary, results)

	return summary, nil
}

// buildSummaryPrompt constructs the prompt for AI
func (s *SearchSummaryService) buildSummaryPrompt(results []SearchResult, query string) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("用户搜索查询: \"%s\"\n\n", query))
	sb.WriteString(fmt.Sprintf("找到 %d 封相关邮件:\n\n", len(results)))

	// Limit to top 10 results for context
	limit := 10
	if len(results) < limit {
		limit = len(results)
	}

	for i := 0; i < limit; i++ {
		result := results[i]
		sb.WriteString(fmt.Sprintf("%d. 发件人: %s\n", i+1, result.Sender))
		sb.WriteString(fmt.Sprintf("   主题: %s\n", result.Subject))
		sb.WriteString(fmt.Sprintf("   摘要: %s\n", truncateText(result.Snippet, 100)))
		sb.WriteString(fmt.Sprintf("   日期: %s\n\n", result.Date.Format("2006-01-02")))
	}

	if len(results) > limit {
		sb.WriteString(fmt.Sprintf("(还有 %d 封邮件未列出)\n\n", len(results)-limit))
	}

	sb.WriteString("请提供:\n")
	sb.WriteString("1. 一句话总结 (中文, 30-50字)\n")
	sb.WriteString("2. 3-5个关键主题\n")
	sb.WriteString("3. 重要联系人 (最多5人)\n")
	sb.WriteString("4. 是否有紧急邮件\n")
	sb.WriteString("\n请以JSON格式返回，格式如下:\n")
	sb.WriteString("{\n")
	sb.WriteString("  \"summary\": \"总结内容\",\n")
	sb.WriteString("  \"topics\": [\"主题1\", \"主题2\"],\n")
	sb.WriteString("  \"people\": [\"人名1\", \"人名2\"],\n")
	sb.WriteString("  \"urgent_count\": 0\n")
	sb.WriteString("}")

	return sb.String()
}

// parseSummaryResponse parses AI response into structured summary
func (s *SearchSummaryService) parseSummaryResponse(response string, results []SearchResult) *SearchResultsSummary {
	// Simple parsing logic - in production, use proper JSON parsing
	summary := &SearchResultsSummary{
		NaturalSummary: response,
	}

	// Extract basic statistics from results
	senderMap := make(map[string]int)
	subjectWords := make(map[string]int)

	for _, result := range results {
		// Count senders
		if result.Sender != "" {
			senderMap[result.Sender]++
		}

		// Extract subject keywords
		words := strings.Fields(strings.ToLower(result.Subject))
		for _, word := range words {
			if len(word) > 3 {
				subjectWords[word]++
			}
		}
	}

	// Get top senders as important people
	type kv struct {
		Key   string
		Value int
	}
	var senderList []kv
	for k, v := range senderMap {
		senderList = append(senderList, kv{k, v})
	}
	// Sort by count
	for i := 0; i < len(senderList)-1; i++ {
		for j := i + 1; j < len(senderList); j++ {
			if senderList[j].Value > senderList[i].Value {
				senderList[i], senderList[j] = senderList[j], senderList[i]
			}
		}
	}

	// Get top 5 senders
	for i := 0; i < len(senderList) && i < 5; i++ {
		summary.ImportantPeople = append(summary.ImportantPeople, senderList[i].Key)
	}

	// Get top keywords as topics
	var wordList []kv
	for k, v := range subjectWords {
		if v > 1 { // Only words appearing more than once
			wordList = append(wordList, kv{k, v})
		}
	}
	// Sort by count
	for i := 0; i < len(wordList)-1; i++ {
		for j := i + 1; j < len(wordList); j++ {
			if wordList[j].Value > wordList[i].Value {
				wordList[i], wordList[j] = wordList[j], wordList[i]
			}
		}
	}

	// Get top 5 topics
	for i := 0; i < len(wordList) && i < 5; i++ {
		summary.KeyTopics = append(summary.KeyTopics, wordList[i].Key)
	}

	// If AI didn't provide a good summary, create a basic one
	if len(response) < 20 {
		summary.NaturalSummary = fmt.Sprintf(
			"找到 %d 封相关邮件，主要来自 %d 位联系人，涉及多个主题。",
			len(results),
			len(senderMap),
		)
	}

	return summary
}

// truncateText truncates text to specified length
func truncateText(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}
	return text[:maxLen] + "..."
}

// GenerateQuickSummary generates a quick non-AI summary
func (s *SearchSummaryService) GenerateQuickSummary(results []SearchResult) *SearchResultsSummary {
	if len(results) == 0 {
		return &SearchResultsSummary{
			NaturalSummary: "未找到相关邮件。",
		}
	}

	// Count unique senders
	senderMap := make(map[string]bool)
	for _, result := range results {
		if result.Sender != "" {
			senderMap[result.Sender] = true
		}
	}

	summary := &SearchResultsSummary{
		NaturalSummary: fmt.Sprintf(
			"找到 %d 封相关邮件，来自 %d 位不同的发件人。",
			len(results),
			len(senderMap),
		),
		ImportantPeople: make([]string, 0),
		KeyTopics:       make([]string, 0),
	}

	// Add top senders
	count := 0
	for sender := range senderMap {
		summary.ImportantPeople = append(summary.ImportantPeople, sender)
		count++
		if count >= 3 {
			break
		}
	}

	return summary
}
