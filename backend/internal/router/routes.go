package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hrygo/echomind/internal/handler"
)

// Handlers holds all HTTP handlers
type Handlers struct {
	Health      *handler.HealthHandler
	Auth        *handler.AuthHandler
	Org         *handler.OrganizationHandler
	Account     *handler.AccountHandler
	Sync        *handler.SyncHandler
	Email       *handler.EmailHandler
	Insight     *handler.InsightHandler
	AIDraft     *handler.AIDraftHandler
	Search      *handler.SearchHandler
	Chat        *handler.ChatHandler
	Task        *handler.TaskHandler
	Context     *handler.ContextHandler
	Action      *handler.ActionHandler
	Opportunity *handler.OpportunityHandler
	WeChat      interface{ Callback(c *gin.Context) } // WeChat gateway handler
}

// SetupRoutes registers all API routes
func SetupRoutes(router *gin.Engine, h *Handlers, authMiddleware gin.HandlerFunc) {
	api := router.Group("/api/v1")
	{
		// Public routes
		api.GET("/health", h.Health.HealthCheck)
		api.POST("/auth/register", h.Auth.Register)
		api.POST("/auth/login", h.Auth.Login)

		// WeChat callback (public)
		if h.WeChat != nil {
			api.Any("/wechat/callback", h.WeChat.Callback)
		}

		// Protected routes
		protected := api.Group("/").Use(authMiddleware)
		{
			// Users
			protected.PATCH("/users/me", h.Auth.UpdateUserProfile)

			// Organization
			protected.POST("/orgs", h.Org.CreateOrganization)
			protected.GET("/orgs", h.Org.ListOrganizations)
			protected.GET("/orgs/:id", h.Org.GetOrganization)
			protected.GET("/orgs/:id/members", h.Org.GetMembers)

			// Account & Sync
			protected.POST("/settings/account", h.Account.ConnectAndSaveAccount)
			protected.GET("/settings/account", h.Account.GetAccountStatus)
			protected.DELETE("/settings/account", h.Account.DisconnectAccount)
			protected.POST("/sync", h.Sync.SyncEmails)

			// Emails & Insights
			protected.GET("/emails", h.Email.ListEmails)
			protected.GET("/emails/:id", h.Email.GetEmail)
			protected.DELETE("/emails/all", h.Email.DeleteAllEmails)
			protected.GET("/insights/network", h.Insight.GetNetworkGraph)

			// Dashboard Insights
			protected.GET("/insights/manager/stats", h.Insight.GetManagerStats)
			protected.GET("/insights/executive/overview", h.Insight.GetExecutiveOverview)
			protected.GET("/insights/risks", h.Insight.GetRisks)
			protected.GET("/insights/trends", h.Insight.GetTrends)
			protected.GET("/insights/dealmaker/radar", h.Insight.GetDealmakerRadar)

			// AI & Search
			protected.POST("/ai/draft", h.AIDraft.GenerateDraft)
			protected.POST("/ai/reply", h.AIDraft.GenerateReply)
			protected.GET("/search", h.Search.Search)
			protected.POST("/chat/completions", h.Chat.StreamChat)

			// Tasks
			protected.POST("/tasks", h.Task.CreateTask)
			protected.GET("/tasks", h.Task.ListTasks)
			protected.PATCH("/tasks/:id", h.Task.UpdateTask)
			protected.PATCH("/tasks/:id/status", h.Task.UpdateTaskStatus)
			protected.DELETE("/tasks/:id", h.Task.DeleteTask)

			// Contexts
			protected.POST("/contexts", h.Context.CreateContext)
			protected.GET("/contexts", h.Context.ListContexts)
			protected.PATCH("/contexts/:id", h.Context.UpdateContext)
			protected.DELETE("/contexts/:id", h.Context.DeleteContext)

			// Actions
			protected.POST("/actions/approve", h.Action.ApproveEmail)
			protected.POST("/actions/snooze", h.Action.SnoozeEmail)
			protected.POST("/actions/dismiss", h.Action.DismissEmail)

			// Opportunities
			protected.POST("/opportunities", h.Opportunity.CreateOpportunity)
			protected.GET("/opportunities", h.Opportunity.ListOpportunities)
			protected.GET("/opportunities/:id", h.Opportunity.GetOpportunity)
			protected.PATCH("/opportunities/:id", h.Opportunity.UpdateOpportunity)
			protected.DELETE("/opportunities/:id", h.Opportunity.DeleteOpportunity)
		}
	}
}
