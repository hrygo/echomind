package model

// ContextInput defines the input structure for creating or updating a context.
type ContextInput struct {
	Name         string   `json:"name" binding:"required,max=100"`
	Color        string   `json:"color" binding:"max=20"`
	Keywords     []string `json:"keywords"`
	Stakeholders []string `json:"stakeholders"` // Email addresses
}
