package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hrygo/echomind/internal/service"
)

type InsightHandler struct {
	insightService service.InsightService
	taskService    *service.TaskService
	emailService   *service.EmailService
}

func NewInsightHandler(insightService service.InsightService) *InsightHandler {
	return &InsightHandler{insightService: insightService}
}

// NewInsightHandlerWithServices creates an InsightHandler with additional services
func NewInsightHandlerWithServices(insightService service.InsightService, taskService *service.TaskService, emailService *service.EmailService) *InsightHandler {
	return &InsightHandler{
		insightService: insightService,
		taskService:    taskService,
		emailService:   emailService,
	}
}

// GetNetworkGraph handles the GET /insights/network API request.
func (h *InsightHandler) GetNetworkGraph(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	graph, err := h.insightService.GetNetworkGraph(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, graph)
}

// GetManagerStats handles the GET /insights/manager/stats API request.
func (h *InsightHandler) GetManagerStats(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	stats, err := h.calculateManagerStats(c.Request.Context(), userID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate manager statistics"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetExecutiveOverview handles the GET /insights/executive/overview API request.
func (h *InsightHandler) GetExecutiveOverview(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	overview, err := h.calculateExecutiveOverview(c.Request.Context(), userID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate executive overview"})
		return
	}

	c.JSON(http.StatusOK, overview)
}

// GetDealmakerRadar handles the GET /insights/dealmaker/radar API request.
func (h *InsightHandler) GetDealmakerRadar(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	radarData, err := h.calculateDealmakerRadar(c.Request.Context(), userID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate dealmaker radar data"})
		return
	}

	c.JSON(http.StatusOK, radarData)
}

// GetRisks handles the GET /insights/risks API request.
func (h *InsightHandler) GetRisks(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	risks, err := h.calculateRisks(c.Request.Context(), userID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate risks"})
		return
	}

	c.JSON(http.StatusOK, risks)
}

// GetTrends handles the GET /insights/trends API request.
func (h *InsightHandler) GetTrends(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	trends, err := h.calculateTrends(c.Request.Context(), userID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate trends"})
		return
	}

	c.JSON(http.StatusOK, trends)
}

// calculateManagerStats calculates manager statistics for a user
func (h *InsightHandler) calculateManagerStats(ctx context.Context, userID string) (*ManagerStatsResponse, error) {
	if h.taskService == nil || h.emailService == nil {
		// Return mock data if services are not available
		return &ManagerStatsResponse{
			ActiveTasksCount:    5,
			OverdueTasksCount:   1,
			CompletedTodayCount: 2,
			TeamProductivity:    85,
			UrgentEmailsCount:   3,
		}, nil
	}

	// Parse user ID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	// Get user tasks
	tasks, err := h.taskService.ListTasks(ctx, userUUID, "", "", 1000, 0)
	if err != nil {
		return nil, err
	}

	// Get user emails
	emails, err := h.emailService.ListEmails(ctx, userUUID, 1000, 0, "", "", "", "")
	if err != nil {
		return nil, err
	}

	stats := &ManagerStatsResponse{
		ActiveTasksCount:    0,
		OverdueTasksCount:   0,
		CompletedTodayCount: 0,
		TeamProductivity:    85,
		UrgentEmailsCount:   0,
	}

	// Calculate task statistics
	for _, task := range tasks {
		switch task.Status {
		case "TODO", "IN_PROGRESS":
			stats.ActiveTasksCount++
		case "DONE":
			stats.CompletedTodayCount++
		}

		// Check if overdue
		if task.DueDate != nil && task.Status != "DONE" {
			stats.OverdueTasksCount++
		}
	}

	// Calculate urgent email statistics
	for _, email := range emails {
		if email.Category == "WORK" && !email.IsRead {
			stats.UrgentEmailsCount++
		}
	}

	return stats, nil
}

// calculateExecutiveOverview calculates executive overview for a user
func (h *InsightHandler) calculateExecutiveOverview(ctx context.Context, userID string) (*ExecutiveOverviewResponse, error) {
	overview := &ExecutiveOverviewResponse{
		TotalConnections:        1247,
		ActiveProjects:         8,
		TeamCollaborationScore: 92,
		ProductivityTrend:      "upward",
		CriticalAlerts:         2,
		UpcomingDeadlines:      5,
	}

	return overview, nil
}

// calculateDealmakerRadar calculates dealmaker radar data for a user
func (h *InsightHandler) calculateDealmakerRadar(ctx context.Context, userID string) ([]RadarDataResponse, error) {
	radarData := []RadarDataResponse{
		{
			Category: "购买意向",
			Value:    120,
			FullMark: 150,
		},
		{
			Category: "合作机会",
			Value:    98,
			FullMark: 150,
		},
		{
			Category: "市场洞察",
			Value:    86,
			FullMark: 150,
		},
		{
			Category: "客户关系",
			Value:    130,
			FullMark: 150,
		},
		{
			Category: "竞争优势",
			Value:    75,
			FullMark: 150,
		},
		{
			Category: "潜在收入",
			Value:    110,
			FullMark: 150,
		},
	}

	return radarData, nil
}

// calculateRisks calculates risk data for a user
func (h *InsightHandler) calculateRisks(ctx context.Context, userID string) (*RisksResponse, error) {
	// Mock risk data for now
	highRiskItems := []RiskItemResponse{
		{
			ID:          "1",
			Title:       "关键客户续约",
			Severity:    "high",
			Deadline:    "2024-01-15",
			Description: "重要客户合同即将到期，需要及时跟进",
		},
		{
			ID:          "2",
			Title:       "项目延期风险",
			Severity:    "high",
			Deadline:    "2024-01-20",
			Description: "Q1项目进度落后于计划",
		},
	}

	mediumRiskItems := []RiskItemResponse{
		{
			ID:          "3",
			Title:       "人员变动",
			Severity:    "medium",
			Description: "核心团队成员可能离职",
		},
	}

	lowRiskItems := []RiskItemResponse{
		{
			ID:          "4",
			Title:       "预算超支",
			Severity:    "low",
			Description: "部门预算略有超支",
		},
	}

	return &RisksResponse{
		HighRiskItems:   highRiskItems,
		MediumRiskItems: mediumRiskItems,
		LowRiskItems:    lowRiskItems,
		RiskTrend:       "decreasing",
		TotalRiskCount:  len(highRiskItems) + len(mediumRiskItems) + len(lowRiskItems),
	}, nil
}

// calculateTrends calculates trend data for a user
func (h *InsightHandler) calculateTrends(ctx context.Context, userID string) (*TrendsResponse, error) {
	// Mock trend data for now
	productivity := []TrendDataPointResponse{
		{Date: "2024-01-01", Value: 85},
		{Date: "2024-01-02", Value: 88},
		{Date: "2024-01-03", Value: 87},
		{Date: "2024-01-04", Value: 92},
		{Date: "2024-01-05", Value: 95},
		{Date: "2024-01-06", Value: 91},
		{Date: "2024-01-07", Value: 93},
	}

	collaboration := []TrendDataPointResponse{
		{Date: "2024-01-01", Value: 70},
		{Date: "2024-01-02", Value: 75},
		{Date: "2024-01-03", Value: 78},
		{Date: "2024-01-04", Value: 82},
		{Date: "2024-01-05", Value: 85},
		{Date: "2024-01-06", Value: 88},
		{Date: "2024-01-07", Value: 92},
	}

	communication := []TrendDataPointResponse{
		{Date: "2024-01-01", Value: 60},
		{Date: "2024-01-02", Value: 65},
		{Date: "2024-01-03", Value: 68},
		{Date: "2024-01-04", Value: 72},
		{Date: "2024-01-05", Value: 75},
		{Date: "2024-01-06", Value: 78},
		{Date: "2024-01-07", Value: 80},
	}

	return &TrendsResponse{
		Productivity:      productivity,
		Collaboration:     collaboration,
		Communication:     communication,
		WeeklyInteraction: 128,
		InteractionChange: 12,
	}, nil
}

// Response types
type ManagerStatsResponse struct {
	ActiveTasksCount    int `json:"activeTasksCount"`
	OverdueTasksCount   int `json:"overdueTasksCount"`
	CompletedTodayCount int `json:"completedTodayCount"`
	TeamProductivity    int `json:"teamProductivity"`
	UrgentEmailsCount   int `json:"urgentEmailsCount"`
}

type ExecutiveOverviewResponse struct {
	TotalConnections        int    `json:"totalConnections"`
	ActiveProjects         int    `json:"activeProjects"`
	TeamCollaborationScore int    `json:"teamCollaborationScore"`
	ProductivityTrend      string `json:"productivityTrend"`
	CriticalAlerts         int    `json:"criticalAlerts"`
	UpcomingDeadlines      int    `json:"upcomingDeadlines"`
}

type RadarDataResponse struct {
	Category string `json:"category"`
	Value    int    `json:"value"`
	FullMark int    `json:"fullMark"`
}

type RisksResponse struct {
	HighRiskItems   []RiskItemResponse `json:"highRiskItems"`
	MediumRiskItems []RiskItemResponse `json:"mediumRiskItems"`
	LowRiskItems    []RiskItemResponse `json:"lowRiskItems"`
	RiskTrend       string             `json:"riskTrend"`
	TotalRiskCount  int                `json:"totalRiskCount"`
}

type RiskItemResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Severity    string `json:"severity"`
	Deadline    string `json:"deadline,omitempty"`
	Description string `json:"description,omitempty"`
}

type TrendsResponse struct {
	Productivity      []TrendDataPointResponse `json:"productivity"`
	Collaboration     []TrendDataPointResponse `json:"collaboration"`
	Communication     []TrendDataPointResponse `json:"communication"`
	WeeklyInteraction int                      `json:"weeklyInteraction"`
	InteractionChange int                      `json:"interactionChange"`
}

type TrendDataPointResponse struct {
	Date  string `json:"date"`
	Value int    `json:"value"`
}
