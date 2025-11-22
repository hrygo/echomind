package model

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TaskStatus string
type TaskPriority string

const (
	TaskStatusTodo       TaskStatus = "todo"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusDone       TaskStatus = "done"

	TaskPriorityHigh   TaskPriority = "high"
	TaskPriorityMedium TaskPriority = "medium"
	TaskPriorityLow    TaskPriority = "low"
)

type Task struct {
	ID            uuid.UUID      `gorm:"type:uuid;primary_key"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`

	UserID        uuid.UUID      `gorm:"type:uuid;not null;index"`
	
	// Optional: Link to source email
	SourceEmailID *uuid.UUID     `gorm:"type:uuid"` 
	
	// Optional: Link to a project context (Week 2)
	ContextID     *uuid.UUID     `gorm:"type:uuid"`

	Title         string         `gorm:"not null"`
	Description   string         `gorm:"type:text"`
	
	Status        TaskStatus     `gorm:"type:varchar(20);default:'todo'"`
	Priority      TaskPriority   `gorm:"type:varchar(20);default:'medium'"`
	
	DueDate       *time.Time
	NotifyWeChat  bool           `gorm:"default:false"` // For future use
}
