package service

import (
	"encoding/json"
	"strings"

	"github.com/google/uuid"
	"github.com/hrygo/echomind/internal/model"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ContextService struct {
	db *gorm.DB
}

func NewContextService(db *gorm.DB) *ContextService {
	return &ContextService{db: db}
}

// CreateContext creates a new context for a user.
func (s *ContextService) CreateContext(userID uuid.UUID, input model.ContextInput) (*model.Context, error) {
	keywordsJSON, err := json.Marshal(input.Keywords)
	if err != nil {
		return nil, err
	}
	stakeholdersJSON, err := json.Marshal(input.Stakeholders)
	if err != nil {
		return nil, err
	}

	ctx := &model.Context{
		ID:           uuid.New(),
		UserID:       userID,
		Name:         input.Name,
		Color:        input.Color,
		Keywords:     datatypes.JSON(keywordsJSON),
		Stakeholders: datatypes.JSON(stakeholdersJSON),
	}

	if err := s.db.Create(ctx).Error; err != nil {
		return nil, err
	}

	return ctx, nil
}

// ListContexts returns all contexts for a user.
func (s *ContextService) ListContexts(userID uuid.UUID) ([]model.Context, error) {
	var contexts []model.Context
	if err := s.db.Where("user_id = ?", userID).Order("created_at desc").Find(&contexts).Error; err != nil {
		return nil, err
	}
	return contexts, nil
}

// GetContext returns a specific context by ID and UserID.
func (s *ContextService) GetContext(id uuid.UUID, userID uuid.UUID) (*model.Context, error) {
	var ctx model.Context
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).First(&ctx).Error; err != nil {
		return nil, err
	}
	return &ctx, nil
}

// UpdateContext updates an existing context.
func (s *ContextService) UpdateContext(id uuid.UUID, userID uuid.UUID, input model.ContextInput) (*model.Context, error) {
	var ctx model.Context
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).First(&ctx).Error; err != nil {
		return nil, err
	}

	keywordsJSON, err := json.Marshal(input.Keywords)
	if err != nil {
		return nil, err
	}
	stakeholdersJSON, err := json.Marshal(input.Stakeholders)
	if err != nil {
		return nil, err
	}

	ctx.Name = input.Name
	ctx.Color = input.Color
	ctx.Keywords = datatypes.JSON(keywordsJSON)
	ctx.Stakeholders = datatypes.JSON(stakeholdersJSON)

	if err := s.db.Save(&ctx).Error; err != nil {
		return nil, err
	}

	return &ctx, nil
}

// DeleteContext deletes a context.
func (s *ContextService) DeleteContext(id uuid.UUID, userID uuid.UUID) error {
	result := s.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Context{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// MatchContexts finds matching contexts for an email based on rules.
func (s *ContextService) MatchContexts(email *model.Email) ([]model.Context, error) {
	// 1. Fetch all contexts for the user
	contexts, err := s.ListContexts(email.UserID)
	if err != nil {
		return nil, err
	}

	var matches []model.Context
	for _, ctx := range contexts {
		matched := false

		// Check Keywords
		var keywords []string
		_ = json.Unmarshal(ctx.Keywords, &keywords)
		for _, kw := range keywords {
			if strings.Contains(strings.ToLower(email.Subject), strings.ToLower(kw)) ||
				strings.Contains(strings.ToLower(email.Snippet), strings.ToLower(kw)) {
				matched = true
				break
			}
		}

		if matched {
			matches = append(matches, ctx)
			continue
		}

		// Check Stakeholders
		var stakeholders []string
		_ = json.Unmarshal(ctx.Stakeholders, &stakeholders)
		for _, sh := range stakeholders {
			if strings.EqualFold(email.Sender, sh) {
				matched = true
				break
			}
			// Future: Check To/CC if available in email model as structured data
		}

		if matched {
			matches = append(matches, ctx)
		}
	}

	return matches, nil
}

// AssignContextsToEmail links contexts to an email.
func (s *ContextService) AssignContextsToEmail(emailID uuid.UUID, contextIDs []uuid.UUID) error {
	if len(contextIDs) == 0 {
		return nil
	}
	var emailContexts []model.EmailContext
	for _, cid := range contextIDs {
		emailContexts = append(emailContexts, model.EmailContext{
			EmailID:   emailID,
			ContextID: cid,
		})
	}
	// Use Clause(clause.OnConflict{DoNothing: true}) to avoid duplicates
	return s.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&emailContexts).Error
}
