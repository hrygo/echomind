package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/hrygo/echomind/internal/model"
	"gorm.io/gorm"
)

// Node represents a contact in the network graph.
type Node struct {
	ID    uuid.UUID `json:"id"`
	Label string    `json:"label"` // Display name for the node (e.g., contact's email or name)
	// Add other relevant contact attributes here if needed for visualization
	InteractionCount int     `json:"interactionCount"`
	AvgSentiment     float64 `json:"avgSentiment"`
}

// Link represents a connection between two nodes (contacts).
// For now, this can be a placeholder. More complex logic is needed to define actual links.
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
// Currently, it only populates nodes based on contacts. Links are empty for now.
func (s *DefaultInsightService) GetNetworkGraph(ctx context.Context, userID uuid.UUID) (*NetworkGraph, error) {
	var contacts []model.Contact
	if err := s.db.WithContext(ctx).Where("user_id = ?", userID).Find(&contacts).Error; err != nil {
		return nil, err
	}

	nodes := make([]Node, len(contacts))
	for i, contact := range contacts {
		nodes[i] = Node{
			ID:               contact.ID,
			Label:            contact.Email, // Use email as label for now, can be improved to use Name if available
			InteractionCount: contact.InteractionCount,
			AvgSentiment:     contact.AvgSentiment,
		}
	}

	// For now, links are empty. This can be extended later.
	return &NetworkGraph{
		Nodes: nodes,
		Links: []Link{},
	}, nil
}
