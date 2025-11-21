package model

// EmailAccountInput defines the input structure for connecting or updating an email account.
type EmailAccountInput struct {
	Email         string `json:"email" binding:"required,email"`
	ServerAddress string `json:"server_address" binding:"required"`
	ServerPort    int    `json:"server_port" binding:"required,min=1,max=65535"`
	Username      string `json:"username" binding:"required"`
	Password      string `json:"password" binding:"required"` // Raw password from user
}
