package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hrygo/echomind/internal/model"
	"github.com/hrygo/echomind/pkg/utils"
	"gorm.io/gorm"
)

type OrganizationService struct {
	db *gorm.DB
}

func NewOrganizationService(db *gorm.DB) *OrganizationService {
	return &OrganizationService{db: db}
}

// CreatePersonalOrganization creates a default organization for a user
func (s *OrganizationService) CreatePersonalOrganization(ctx context.Context, user *model.User, tx *gorm.DB) (*model.Organization, error) {
	org := &model.Organization{
		ID:      uuid.New(),
		Name:    fmt.Sprintf("%s's Workspace", user.Name),
		Slug:    fmt.Sprintf("%s-workspace-%s", user.Name, uuid.NewString()[:8]), // Simple unique slug
		OwnerID: user.ID,
	}

	db := s.db
	if tx != nil {
		db = tx
	}

	err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create Organization
		if err := tx.Create(org).Error; err != nil {
			return err
		}

		// Add User as Owner
		member := &model.OrganizationMember{
			OrganizationID: org.ID,
			UserID:         user.ID,
			Role:           model.OrgRoleOwner,
		}
		if err := tx.Create(member).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return org, nil
}

// GetOrganizationMembers returns members of an org

func (s *OrganizationService) GetOrganizationMembers(ctx context.Context, orgID uuid.UUID) ([]model.OrganizationMember, error) {

	var members []model.OrganizationMember

	// Preload User to get names/emails

	err := s.db.WithContext(ctx).
		Preload("User").
		Where("organization_id = ?", orgID).
		Find(&members).Error

	return members, err

}

// AddMember adds a user to an organization

func (s *OrganizationService) AddMember(ctx context.Context, orgID, userID uuid.UUID, role model.OrganizationRole) error {

	member := &model.OrganizationMember{

		OrganizationID: orgID,

		UserID: userID,

		Role: role,
	}

	return s.db.WithContext(ctx).Create(member).Error

}

// EnsureAllUsersHaveOrganization runs a migration to check all users
func (s *OrganizationService) EnsureAllUsersHaveOrganization(ctx context.Context) error {
	var users []model.User
	// Find users who are NOT in any organization
	// Subquery: SELECT user_id FROM organization_members
	err := s.db.WithContext(ctx).
		Where("id NOT IN (?)", s.db.Table("organization_members").Select("user_id")).
		Find(&users).Error

	if err != nil {
		return err
	}

	count := 0
	for _, user := range users {
		if _, err := s.CreatePersonalOrganization(ctx, &user, nil); err != nil {
			// Log error but continue? Or fail hard?
			// For now, let's return error to stop migration if something is wrong
			return fmt.Errorf("failed to create org for user %s: %w", user.ID, err)
		}
		count++
	}

	if count > 0 {
		fmt.Printf("Migrated %d users to personal organizations\n", count)
	}

	return nil
}

// GetUserOrganizations returns all organizations a user belongs to
func (s *OrganizationService) GetUserOrganizations(ctx context.Context, userID uuid.UUID) ([]model.Organization, error) {
	var orgs []model.Organization
	// Join with members table
	err := s.db.WithContext(ctx).
		Joins("JOIN organization_members ON organization_members.organization_id = organizations.id").
		Where("organization_members.user_id = ?", userID).
		Find(&orgs).Error
	return orgs, err
}

// GetOrganizationByID returns org details if user is a member
func (s *OrganizationService) GetOrganizationByID(ctx context.Context, orgID, userID uuid.UUID) (*model.Organization, error) {
	var org model.Organization
	// Verify membership implicitly
	err := s.db.WithContext(ctx).
		Joins("JOIN organization_members ON organization_members.organization_id = organizations.id").
		Where("organizations.id = ? AND organization_members.user_id = ?", orgID, userID).
		First(&org).Error

	if err != nil {
		return nil, err
	}
	return &org, nil
}

// CreateOrganization creates a new organization
func (s *OrganizationService) CreateOrganization(ctx context.Context, name string, ownerID uuid.UUID) (*model.Organization, error) {
	org := &model.Organization{
		ID:      uuid.New(),
		Name:    name,
		Slug:    fmt.Sprintf("%s-%s", utils.Slugify(name), uuid.NewString()[:8]),
		OwnerID: ownerID,
	}

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(org).Error; err != nil {
			return err
		}
		// Add owner
		member := &model.OrganizationMember{
			OrganizationID: org.ID,
			UserID:         ownerID,
			Role:           model.OrgRoleOwner,
		}
		return tx.Create(member).Error
	})

	if err != nil {
		return nil, err
	}
	return org, nil
}
