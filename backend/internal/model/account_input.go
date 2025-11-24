package model

// EmailAccountInput defines the input structure for connecting or updating an email account.
type EmailAccountInput struct {
	Email          string  `json:"email" binding:"required,email"`
	ServerAddress  string  `json:"server_address" binding:"required"`
	ServerPort     int     `json:"server_port" binding:"required,min=1,max=65535"`
	Username       string  `json:"username" binding:"required"`
	SMTPServer     string  `json:"smtp_server" binding:"required"`
	SMTPPort       int     `json:"smtp_port" binding:"required,min=1,max=65535"`
	Password       string  `json:"password" binding:"required"` // Raw password from user
	TeamID         *string `json:"team_id"`                     // Optional, UUID as string
	OrganizationID *string `json:"organization_id"`             // Optional, UUID as string
}
