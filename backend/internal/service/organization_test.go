package service

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/hrygo/echomind/internal/model"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	_ = db.AutoMigrate(&model.User{}, &model.Organization{}, &model.OrganizationMember{})
	return db
}

func TestCreatePersonalOrganization(t *testing.T) {
	db := setupTestDB()
	svc := NewOrganizationService(db)
	ctx := context.Background()

	user := &model.User{
		ID:    uuid.New(),
		Name:  "Test User",
		Email: "test@example.com",
	}
	// Create Personal Org
	org, err := svc.CreatePersonalOrganization(ctx, user, nil)
	assert.NoError(t, err)
	assert.NotNil(t, org)
	assert.Equal(t, "Test User's Workspace", org.Name)
	assert.Equal(t, user.ID, org.OwnerID)

	// Verify Membership
	var member model.OrganizationMember
	err = db.Where("organization_id = ? AND user_id = ?", org.ID, user.ID).First(&member).Error
	assert.NoError(t, err)
	assert.Equal(t, model.OrgRoleOwner, member.Role)
}

func TestEnsureAllUsersHaveOrganization(t *testing.T) {
	db := setupTestDB()
	svc := NewOrganizationService(db)
	ctx := context.Background()

	// Create 2 users without orgs
	u1 := &model.User{ID: uuid.New(), Name: "User1", Email: "u1@example.com"}
	u2 := &model.User{ID: uuid.New(), Name: "User2", Email: "u2@example.com"}
	db.Create(u1)
	db.Create(u2)

	err := svc.EnsureAllUsersHaveOrganization(ctx)
	assert.NoError(t, err)

	// Verify orgs created
	var count int64
	db.Model(&model.Organization{}).Count(&count)
	assert.Equal(t, int64(2), count)

	// Run again (idempotency)
	err = svc.EnsureAllUsersHaveOrganization(ctx)
	assert.NoError(t, err)
	db.Model(&model.Organization{}).Count(&count)
	assert.Equal(t, int64(2), count)
}
