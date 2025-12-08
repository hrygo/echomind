package service

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hrygo/echomind/internal/model"
	"github.com/stretchr/testify/assert"
	"gorm.io/datatypes"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupInsightTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	// Migrate schema
	err = db.AutoMigrate(&model.Contact{}, &model.Email{})
	if err != nil {
		t.Fatalf("Failed to migrate schema: %v", err)
	}

	return db
}

func TestGetNetworkGraph(t *testing.T) {
	db := setupInsightTestDB(t)
	svc := NewInsightService(db)
	ctx := context.Background()

	// 1. Setup Data
	userID := uuid.New()

	// Create Contacts
	alice := model.Contact{
		ID:     uuid.New(),
		UserID: &userID,
		Email:  "alice@example.com",
		Name:   "Alice",
	}
	bob := model.Contact{
		ID:     uuid.New(),
		UserID: &userID,
		Email:  "bob@example.com",
		Name:   "Bob",
	}
	charlie := model.Contact{
		ID:     uuid.New(),
		UserID: &userID,
		Email:  "charlie@example.com",
		Name:   "Charlie",
	}
	// Unknown person (not in contacts)
	unknownEmail := "unknown@example.com"

	db.Create(&alice)
	db.Create(&bob)
	db.Create(&charlie)

	// Create Emails
	// Email 1: Alice -> Bob (Link: Alice-Bob)
	email1 := model.Email{
		ID:        uuid.New(),
		UserID:    userID,
		Sender:    alice.Email,
		To:        datatypes.JSON(jsonRaw([]string{bob.Email})),
		Cc:        datatypes.JSON(jsonRaw([]string{})),
		Date:      time.Now(),
		MessageID: "msg1",
	}

	// Email 2: Bob -> Alice, Charlie (Links: Bob-Alice, Bob-Charlie, Alice-Charlie)
	email2 := model.Email{
		ID:        uuid.New(),
		UserID:    userID,
		Sender:    bob.Email,
		To:        datatypes.JSON(jsonRaw([]string{alice.Email, charlie.Email})),
		Cc:        datatypes.JSON(jsonRaw([]string{})),
		Date:      time.Now(),
		MessageID: "msg2",
	}

	// Email 3: Alice -> Unknown (No links should be created for Unknown)
	email3 := model.Email{
		ID:        uuid.New(),
		UserID:    userID,
		Sender:    alice.Email,
		To:        datatypes.JSON(jsonRaw([]string{unknownEmail})),
		Cc:        datatypes.JSON(jsonRaw([]string{})),
		Date:      time.Now(),
		MessageID: "msg3",
	}

	db.Create(&email1)
	db.Create(&email2)
	db.Create(&email3)

	// 2. Execute
	graph, err := svc.GetNetworkGraph(ctx, userID)

	// 3. Verify
	assert.NoError(t, err)
	assert.NotNil(t, graph)

	// Verify Nodes
	assert.Equal(t, 3, len(graph.Nodes), "Should have 3 nodes (Alice, Bob, Charlie)")

	// Verify Links
	// Expected Links:
	// 1. Alice-Bob (from Email 1)
	// 2. Bob-Alice (from Email 2) -> Total Weight 2
	// 3. Bob-Charlie (from Email 2) -> Total Weight 1
	// 4. Alice-Charlie (from Email 2, co-recipients) -> Total Weight 1
	assert.Equal(t, 3, len(graph.Links), "Should have 3 unique links")

	// Helper to find link weight
	getWeight := func(id1, id2 uuid.UUID) float64 {
		for _, l := range graph.Links {
			if (l.Source == id1 && l.Target == id2) || (l.Source == id2 && l.Target == id1) {
				return l.Weight
			}
		}
		return 0
	}

	weightAliceBob := getWeight(alice.ID, bob.ID)
	assert.Equal(t, 2.0, weightAliceBob, "Alice-Bob weight should be 2")

	weightBobCharlie := getWeight(bob.ID, charlie.ID)
	assert.Equal(t, 1.0, weightBobCharlie, "Bob-Charlie weight should be 1")

	weightAliceCharlie := getWeight(alice.ID, charlie.ID)
	assert.Equal(t, 1.0, weightAliceCharlie, "Alice-Charlie weight should be 1")
}

func jsonRaw(v interface{}) []byte {
	b, _ := json.Marshal(v)
	return b
}
