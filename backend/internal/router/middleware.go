package router

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/hrygo/echomind/configs"
	"github.com/hrygo/echomind/internal/bootstrap"
	"github.com/hrygo/echomind/internal/middleware"
	echologger "github.com/hrygo/echomind/pkg/logger"
)

// SetupMiddleware configures all global middleware for the Gin engine
func SetupMiddleware(r *gin.Engine, app *bootstrap.App, isProduction bool) {
	// Configure Gin mode
	if isProduction {
		gin.SetMode(gin.ReleaseMode)
	}

	// Configure trusted proxies (security best practice)
	if err := r.SetTrustedProxies([]string{"127.0.0.1", "::1"}); err != nil {
		app.Logger.Warn("Failed to set trusted proxies", echologger.Error(err))
	}

	// Middleware: Request ID (for tracing)
	r.Use(middleware.RequestID())

	// Middleware: CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Request-ID"},
		ExposeHeaders:    []string{"Content-Length", "X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
}

// SetupAuthMiddleware returns the authentication middleware
func SetupAuthMiddleware(jwtConfig configs.JWTConfig) gin.HandlerFunc {
	return middleware.AuthMiddleware(jwtConfig)
}
