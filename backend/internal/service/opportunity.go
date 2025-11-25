package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/hrygo/echomind/internal/model"
)

// OpportunityService handles opportunity-related business logic
type OpportunityService struct {
	db *gorm.DB
}

// NewOpportunityService creates a new OpportunityService
func NewOpportunityService(db *gorm.DB) *OpportunityService {
	return &OpportunityService{db: db}
}

// CreateOpportunity creates a new opportunity
func (s *OpportunityService) CreateOpportunity(ctx context.Context, userID, teamID, orgID string, title, description, company, value string, oppType model.OpportunityType, confidence int, sourceEmailID *string) (*model.Opportunity, error) {
	opportunity := &model.Opportunity{
		ID:            uuid.New().String(),
		Title:         title,
		Description:   description,
		Company:       company,
		Value:         value,
		Type:          oppType,
		Status:        model.OpportunityStatusNew,
		Confidence:    confidence,
		UserID:        userID,
		TeamID:        teamID,
		OrgID:         orgID,
		SourceEmailID: sourceEmailID,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.db.WithContext(ctx).Create(opportunity).Error; err != nil {
		return nil, fmt.Errorf("failed to create opportunity: %w", err)
	}

	return opportunity, nil
}

// ListOpportunities retrieves opportunities for a user with optional filters
func (s *OpportunityService) ListOpportunities(ctx context.Context, userID string, status *model.OpportunityStatus, oppType *model.OpportunityType, limit, offset int) ([]model.Opportunity, error) {
	var opportunities []model.Opportunity
	query := s.db.WithContext(ctx).Where("user_id = ?", userID)

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if oppType != nil {
		query = query.Where("type = ?", *oppType)
	}

	if limit > 0 {
		query = query.Limit(limit)
	}

	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Order("created_at DESC").Find(&opportunities).Error
	if err != nil {
		return nil, fmt.Errorf("failed to list opportunities: %w", err)
	}

	return opportunities, nil
}

// GetOpportunityByID retrieves an opportunity by ID
func (s *OpportunityService) GetOpportunityByID(ctx context.Context, userID, id string) (*model.Opportunity, error) {
	var opportunity model.Opportunity
	err := s.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&opportunity).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("opportunity not found")
		}
		return nil, fmt.Errorf("failed to get opportunity: %w", err)
	}

	return &opportunity, nil
}

// UpdateOpportunity updates an existing opportunity
func (s *OpportunityService) UpdateOpportunity(ctx context.Context, userID, id string, updates map[string]interface{}) (*model.Opportunity, error) {
	var opportunity model.Opportunity

	// First check if the opportunity exists and belongs to the user
	err := s.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&opportunity).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("opportunity not found")
		}
		return nil, fmt.Errorf("failed to get opportunity: %w", err)
	}

	// Apply updates
	if err := s.db.WithContext(ctx).Model(&opportunity).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("failed to update opportunity: %w", err)
	}

	// Refresh the data
	if err := s.db.WithContext(ctx).Where("id = ?", id).First(&opportunity).Error; err != nil {
		return nil, fmt.Errorf("failed to refresh opportunity data: %w", err)
	}

	return &opportunity, nil
}

// DeleteOpportunity deletes an opportunity
func (s *OpportunityService) DeleteOpportunity(ctx context.Context, userID, id string) error {
	result := s.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&model.Opportunity{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete opportunity: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("opportunity not found")
	}

	return nil
}

// GetOpportunityStats returns statistics about opportunities for a user
func (s *OpportunityService) GetOpportunityStats(ctx context.Context, userID string) (map[string]int, error) {
	var stats []struct {
		Status string
		Count  int
	}

	err := s.db.WithContext(ctx).
		Table("opportunities").
		Select("status, COUNT(*) as count").
		Where("user_id = ?", userID).
		Group("status").
		Find(&stats).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get opportunity stats: %w", err)
	}

	result := make(map[string]int)
	for _, stat := range stats {
		result[stat.Status] = stat.Count
	}

	return result, nil
}