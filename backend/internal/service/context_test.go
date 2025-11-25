package service

import (
	"testing"

	"github.com/google/uuid"
	"github.com/hrygo/echomind/internal/model"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupContextTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	_ = db.AutoMigrate(&model.Context{}, &model.EmailContext{}, &model.Email{})
	return db
}

func TestMatchContexts(t *testing.T) {
	db := setupContextTestDB()
	svc := NewContextService(db)
	userID := uuid.New()

	// 1. Create Contexts
	kw := []string{"Project Alpha", "Q4 Budget"}
	sh := []string{"boss@example.com"}

	input := model.ContextInput{
		Name:         "Important Project",
		Keywords:     kw,
		Stakeholders: sh,
	}
	_, err := svc.CreateContext(userID, input)
	assert.NoError(t, err)

	// 2. Test Match by Keyword (Subject)
	email1 := &model.Email{
		ID:       uuid.New(),
		UserID:   userID,
		Subject:  "Update on Project Alpha",
		BodyText: "Things are going well.",
		Sender:   "other@example.com",
	}
	matches1, err := svc.MatchContexts(email1)
	assert.NoError(t, err)
	assert.Len(t, matches1, 1)
	assert.Equal(t, "Important Project", matches1[0].Name)

	// 3. Test Match by Stakeholder
	email2 := &model.Email{
		ID:      uuid.New(),
		UserID:  userID,
		Subject: "Lunch?",
		Sender:  "boss@example.com",
	}
	matches2, err := svc.MatchContexts(email2)
	assert.NoError(t, err)
	assert.Len(t, matches2, 1)

	// 4. Test No Match
	email3 := &model.Email{
		ID:      uuid.New(),
		UserID:  userID,
		Subject: "Random spam",
		Sender:  "spammer@example.com",
	}
	matches3, err := svc.MatchContexts(email3)
	assert.NoError(t, err)
	assert.Len(t, matches3, 0)
}

func TestCRUDContext(t *testing.T) {
	db := setupContextTestDB()
	svc := NewContextService(db)
	userID := uuid.New()

	input := model.ContextInput{
		Name: "Test Context",
	}

	// Create
	ctx, err := svc.CreateContext(userID, input)
	assert.NoError(t, err)
	assert.NotNil(t, ctx)
	assert.Equal(t, "Test Context", ctx.Name)

	// List
	list, err := svc.ListContexts(userID)
	assert.NoError(t, err)
	assert.Len(t, list, 1)

	// Update
	input.Name = "Updated Context"
	updated, err := svc.UpdateContext(ctx.ID, userID, input)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Context", updated.Name)

	// Delete
	err = svc.DeleteContext(ctx.ID, userID)
	assert.NoError(t, err)

	// List again
	list, err = svc.ListContexts(userID)
	assert.NoError(t, err)
	assert.Len(t, list, 0)
}
