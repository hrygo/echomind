package model

import (
	"time"

	"gorm.io/gorm"
)

type Contact struct {
	gorm.Model
	Email            string `gorm:"uniqueIndex;not null"`
	Name             string
	InteractionCount int
	LastInteractedAt time.Time
}
