package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HealthHandler struct {
	db *gorm.DB
}

func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

func (h *HealthHandler) HealthCheck(c *gin.Context) {
	// Check Database Connection
	sqlDB, err := h.db.DB()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "error", "details": "database connection unavailable"})
		return
	}

	if err := sqlDB.Ping(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "error", "details": "database ping failed"})
		return
	}

	// Check Vector Extension (optional, but good for this app)
	var version string
	if err := h.db.Raw("SELECT extversion FROM pg_extension WHERE extname = 'vector'").Scan(&version).Error; err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "warning", "details": "vector extension missing"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"services": gin.H{
			"database": "connected",
			"pgvector": version,
		},
	})
}
