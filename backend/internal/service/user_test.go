package service

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hrygo/echomind/internal/model"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupUserTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	_ = db.AutoMigrate(&model.User{})
	return db
}

func TestUpdateUserProfile(t *testing.T) {
	db := setupUserTestDB()
	svc := &UserService{db: db} // Simplified initialization for this test
	ctx := context.Background()

	// 1. Create User
	user := &model.User{
		ID:        uuid.New(),
		Email:     "test@example.com",
		Name:      "Old Name",
		Role:      "executive",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	db.Create(user)

	// 2. Update Name only
	err := svc.UpdateUserProfile(ctx, user.ID, "", "New Name")
	assert.NoError(t, err)

	var updatedUser model.User
	db.First(&updatedUser, "id = ?", user.ID)
	assert.Equal(t, "New Name", updatedUser.Name)
	assert.Equal(t, "executive", updatedUser.Role) // Role unchanged

	// 3. Update Role only
	err = svc.UpdateUserProfile(ctx, user.ID, "manager", "")
	assert.NoError(t, err)

	db.First(&updatedUser, "id = ?", user.ID)
	assert.Equal(t, "New Name", updatedUser.Name) // Name unchanged
	assert.Equal(t, "manager", updatedUser.Role)

	// 4. Update Both
	err = svc.UpdateUserProfile(ctx, user.ID, "dealmaker", "Final Name")
	assert.NoError(t, err)

	db.First(&updatedUser, "id = ?", user.ID)
	assert.Equal(t, "Final Name", updatedUser.Name)
	assert.Equal(t, "dealmaker", updatedUser.Role)
}
