package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hrygo/echomind/internal/model"
	"gorm.io/gorm"
)

// Node represents a contact in the network graph.
type Node struct {
	ID               uuid.UUID `json:"id"`
	Label            string    `json:"label"` // Display name for the node (e.g., contact's email or name)
	InteractionCount int       `json:"interactionCount"`
	AvgSentiment     float64   `json:"avgSentiment"`
}

// Link represents a connection between two nodes (contacts).
type Link struct {
	Source uuid.UUID `json:"source"`
	Target uuid.UUID `json:"target"`
	Weight float64   `json:"weight"` // Strength of the relationship
}

// NetworkGraph holds the nodes and links for the relationship network.
type NetworkGraph struct {
	Nodes []Node `json:"nodes"`
	Links []Link `json:"links"`
}

// InsightService defines the interface for retrieving insights data.
type InsightService interface {
	GetNetworkGraph(ctx context.Context, userID uuid.UUID) (*NetworkGraph, error)
}

// DefaultInsightService is the default implementation of InsightService.
type DefaultInsightService struct {
	db *gorm.DB
}

// NewInsightService creates a new DefaultInsightService.
func NewInsightService(db *gorm.DB) *DefaultInsightService {
	return &DefaultInsightService{db: db}
}

// GetNetworkGraph retrieves all contacts for a user and formats them as a network graph.
// It builds links based on co-occurrence in email threads (Sender/To/Cc).
func (s *DefaultInsightService) GetNetworkGraph(ctx context.Context, userID uuid.UUID) (*NetworkGraph, error) {
	// 1. Fetch Contacts (Nodes)
	var contacts []model.Contact
	if err := s.db.WithContext(ctx).Where("user_id = ?", userID).Find(&contacts).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch contacts: %w", err)
	}

	// Map Email -> ContactID for quick lookup
	emailToID := make(map[string]uuid.UUID)
	nodes := make([]Node, len(contacts))
	for i, contact := range contacts {
		emailToID[contact.Email] = contact.ID
		nodes[i] = Node{
			ID:               contact.ID,
			Label:            contact.Email, // Improve to use Name if available
			InteractionCount: contact.InteractionCount,
			AvgSentiment:     contact.AvgSentiment,
		}
	}

	// 2. Build Links from Emails
	// Fetch recent emails to build the graph dynamically
	var emails []model.Email
	if err := s.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("date DESC").
		Limit(500). // Limit to recent 500 emails for performance
		Find(&emails).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch emails: %w", err)
	}

	linkWeights := make(map[string]float64)

	for _, email := range emails {
		// Collect all participants in this email
		participants := make(map[string]bool)
		participants[email.Sender] = true

		// Helper to unmarshal JSON lists
		var toList []string
		_ = json.Unmarshal(email.To, &toList)
		for _, addr := range toList {
			participants[addr] = true
		}

		var ccList []string
		_ = json.Unmarshal(email.Cc, &ccList)
		for _, addr := range ccList {
			participants[addr] = true
		}

		// Create links between all unique pairs
		participantList := make([]string, 0, len(participants))
		for p := range participants {
			if _, exists := emailToID[p]; exists { // Only link known contacts
				participantList = append(participantList, p)
			}
		}

		for i := 0; i < len(participantList); i++ {
			for j := i + 1; j < len(participantList); j++ {
				p1 := participantList[i]
				p2 := participantList[j]
				id1 := emailToID[p1]
				id2 := emailToID[p2]

				// Ensure consistent key order (Source < Target)
				key := ""
				if id1.String() < id2.String() {
					key = fmt.Sprintf("%s|%s", id1.String(), id2.String())
				} else {
					key = fmt.Sprintf("%s|%s", id2.String(), id1.String())
				}

				linkWeights[key]++
			}
		}
	}

	// Convert map to Links slice
	links := make([]Link, 0, len(linkWeights))
	for key, weight := range linkWeights {
		parts := strings.Split(key, "|")
		if len(parts) != 2 {
			continue
		}

		sourceID, _ := uuid.Parse(parts[0])
		targetID, _ := uuid.Parse(parts[1])

		links = append(links, Link{
			Source: sourceID,
			Target: targetID,
			Weight: weight,
		})
	}

	return &NetworkGraph{
		Nodes: nodes,
		Links: links,
	}, nil
}
