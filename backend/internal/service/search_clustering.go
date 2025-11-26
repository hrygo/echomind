package service

import (
	"sort"
	"strings"
	"time"
)

// ClusterType defines the type of clustering
type ClusterType string

const (
	ClusterByTopic  ClusterType = "topic"
	ClusterBySender ClusterType = "sender"
	ClusterByTime   ClusterType = "time"
)

// SearchCluster represents a cluster of search results
type SearchCluster struct {
	ID      string         `json:"id"`
	Type    ClusterType    `json:"type"`
	Label   string         `json:"label"`
	Count   int            `json:"count"`
	Results []SearchResult `json:"results"`
}

// ClusteredSearchResults contains both raw results and clusters
type ClusteredSearchResults struct {
	Results  []SearchResult  `json:"results"`
	Clusters []SearchCluster `json:"clusters"`
	Summary  ClusterSummary  `json:"summary"`
}

// ClusterSummary provides statistical overview
type ClusterSummary struct {
	TotalResults   int               `json:"total_results"`
	ClusterCount   int               `json:"cluster_count"`
	TopSenders     []SenderStat      `json:"top_senders"`
	TimeRanges     []TimeRangeStat   `json:"time_ranges"`
	TopicKeywords  []string          `json:"topic_keywords,omitempty"`
}

// SenderStat represents sender statistics
type SenderStat struct {
	Sender string `json:"sender"`
	Count  int    `json:"count"`
}

// TimeRangeStat represents time range statistics
type TimeRangeStat struct {
	Label string `json:"label"` // "today", "this_week", "this_month", "older"
	Count int    `json:"count"`
}

// SearchClusteringService handles result clustering
type SearchClusteringService struct{}

// NewSearchClusteringService creates a new clustering service
func NewSearchClusteringService() *SearchClusteringService {
	return &SearchClusteringService{}
}

// ClusterResults clusters search results by specified type
func (s *SearchClusteringService) ClusterResults(results []SearchResult, clusterType ClusterType) []SearchCluster {
	switch clusterType {
	case ClusterBySender:
		return s.clusterBySender(results)
	case ClusterByTime:
		return s.clusterByTime(results)
	case ClusterByTopic:
		return s.clusterByTopic(results)
	default:
		return nil
	}
}

// ClusterAll performs all clustering types and generates summary
func (s *SearchClusteringService) ClusterAll(results []SearchResult) *ClusteredSearchResults {
	return &ClusteredSearchResults{
		Results: results,
		Clusters: []SearchCluster{
			{
				ID:      "sender_clusters",
				Type:    ClusterBySender,
				Label:   "按发件人分组",
				Count:   0,
				Results: results,
			},
			{
				ID:      "time_clusters",
				Type:    ClusterByTime,
				Label:   "按时间分组",
				Count:   0,
				Results: results,
			},
		},
		Summary: s.generateSummary(results),
	}
}

// clusterBySender groups results by sender
func (s *SearchClusteringService) clusterBySender(results []SearchResult) []SearchCluster {
	senderMap := make(map[string][]SearchResult)
	
	for _, result := range results {
		sender := result.Sender
		if sender == "" {
			sender = "Unknown"
		}
		senderMap[sender] = append(senderMap[sender], result)
	}

	clusters := make([]SearchCluster, 0, len(senderMap))
	for sender, senderResults := range senderMap {
		clusters = append(clusters, SearchCluster{
			ID:      generateClusterID("sender", sender),
			Type:    ClusterBySender,
			Label:   sender,
			Count:   len(senderResults),
			Results: senderResults,
		})
	}

	// Sort by count descending
	sort.Slice(clusters, func(i, j int) bool {
		return clusters[i].Count > clusters[j].Count
	})

	return clusters
}

// clusterByTime groups results by time periods
func (s *SearchClusteringService) clusterByTime(results []SearchResult) []SearchCluster {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	thisWeekStart := today.AddDate(0, 0, -int(now.Weekday()))
	thisMonthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	timeGroups := map[string][]SearchResult{
		"today":      {},
		"this_week":  {},
		"this_month": {},
		"older":      {},
	}

	for _, result := range results {
		switch {
		case result.Date.After(today):
			timeGroups["today"] = append(timeGroups["today"], result)
		case result.Date.After(thisWeekStart):
			timeGroups["this_week"] = append(timeGroups["this_week"], result)
		case result.Date.After(thisMonthStart):
			timeGroups["this_month"] = append(timeGroups["this_month"], result)
		default:
			timeGroups["older"] = append(timeGroups["older"], result)
		}
	}

	labels := map[string]string{
		"today":      "今天",
		"this_week":  "本周",
		"this_month": "本月",
		"older":      "更早",
	}

	clusters := make([]SearchCluster, 0, 4)
	for key, timeResults := range timeGroups {
		if len(timeResults) > 0 {
			clusters = append(clusters, SearchCluster{
				ID:      generateClusterID("time", key),
				Type:    ClusterByTime,
				Label:   labels[key],
				Count:   len(timeResults),
				Results: timeResults,
			})
		}
	}

	// Sort by time period (today first)
	order := map[string]int{"today": 0, "this_week": 1, "this_month": 2, "older": 3}
	sort.Slice(clusters, func(i, j int) bool {
		iKey := strings.TrimPrefix(clusters[i].ID, "time_")
		jKey := strings.TrimPrefix(clusters[j].ID, "time_")
		return order[iKey] < order[jKey]
	})

	return clusters
}

// clusterByTopic groups results by simple topic extraction
func (s *SearchClusteringService) clusterByTopic(results []SearchResult) []SearchCluster {
	// Simple keyword-based topic clustering
	// Extract common words from subjects and snippets
	topicMap := make(map[string][]SearchResult)
	
	for _, result := range results {
		// Extract key terms from subject
		keywords := s.extractKeywords(result.Subject + " " + result.Snippet)
		
		if len(keywords) == 0 {
			topicMap["其他"] = append(topicMap["其他"], result)
			continue
		}
		
		// Use first keyword as topic
		topic := keywords[0]
		topicMap[topic] = append(topicMap[topic], result)
	}

	clusters := make([]SearchCluster, 0, len(topicMap))
	for topic, topicResults := range topicMap {
		clusters = append(clusters, SearchCluster{
			ID:      generateClusterID("topic", topic),
			Type:    ClusterByTopic,
			Label:   topic,
			Count:   len(topicResults),
			Results: topicResults,
		})
	}

	// Sort by count descending
	sort.Slice(clusters, func(i, j int) bool {
		return clusters[i].Count > clusters[j].Count
	})

	return clusters
}

// extractKeywords extracts simple keywords from text
func (s *SearchClusteringService) extractKeywords(text string) []string {
	// Convert to lowercase
	text = strings.ToLower(text)
	
	// Split by common separators
	words := strings.FieldsFunc(text, func(r rune) bool {
		return r == ' ' || r == ',' || r == '.' || r == '!' || r == '?' || r == ';' || r == ':'
	})
	
	// Filter out common words and short words
	stopWords := map[string]bool{
		"the": true, "a": true, "an": true, "and": true, "or": true, "but": true,
		"in": true, "on": true, "at": true, "to": true, "for": true, "of": true,
		"with": true, "by": true, "from": true, "is": true, "are": true, "was": true,
		"re": true, "fwd": true, "fw": true, // Email prefixes
	}
	
	keywords := make([]string, 0)
	seen := make(map[string]bool)
	
	for _, word := range words {
		word = strings.TrimSpace(word)
		if len(word) < 3 || stopWords[word] || seen[word] {
			continue
		}
		keywords = append(keywords, word)
		seen[word] = true
		
		if len(keywords) >= 3 {
			break
		}
	}
	
	return keywords
}

// generateSummary creates statistical summary
func (s *SearchClusteringService) generateSummary(results []SearchResult) ClusterSummary {
	// Count by sender
	senderCount := make(map[string]int)
	for _, result := range results {
		sender := result.Sender
		if sender == "" {
			sender = "Unknown"
		}
		senderCount[sender]++
	}

	// Get top senders
	topSenders := make([]SenderStat, 0, len(senderCount))
	for sender, count := range senderCount {
		topSenders = append(topSenders, SenderStat{
			Sender: sender,
			Count:  count,
		})
	}
	sort.Slice(topSenders, func(i, j int) bool {
		return topSenders[i].Count > topSenders[j].Count
	})
	if len(topSenders) > 5 {
		topSenders = topSenders[:5]
	}

	// Count by time ranges
	timeClusters := s.clusterByTime(results)
	timeRanges := make([]TimeRangeStat, 0, len(timeClusters))
	for _, cluster := range timeClusters {
		timeRanges = append(timeRanges, TimeRangeStat{
			Label: cluster.Label,
			Count: cluster.Count,
		})
	}

	return ClusterSummary{
		TotalResults: len(results),
		ClusterCount: len(senderCount) + len(timeClusters),
		TopSenders:   topSenders,
		TimeRanges:   timeRanges,
	}
}

// generateClusterID creates a unique cluster ID
func generateClusterID(prefix, label string) string {
	// Simple ID generation
	cleanLabel := strings.ReplaceAll(strings.ToLower(label), " ", "_")
	return prefix + "_" + cleanLabel
}
