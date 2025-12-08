package handlers

import (
	"net/http"
	"time"

	"github.com/bishworup11/bdSeeker-backend/internal/database"
	"github.com/bishworup11/bdSeeker-backend/internal/middleware"
	"github.com/bishworup11/bdSeeker-backend/internal/models"
	"github.com/bishworup11/bdSeeker-backend/internal/repositories"
	"github.com/bishworup11/bdSeeker-backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	userRepo    *repositories.UserRepository
	companyRepo *repositories.CompanyRepository
	reportRepo  *repositories.ReportRepository
}

func NewAdminHandler() *AdminHandler {
	db := database.GetDB()
	return &AdminHandler{
		userRepo:    repositories.NewUserRepository(db),
		companyRepo: repositories.NewCompanyRepository(db),
		reportRepo:  repositories.NewReportRepository(db),
	}
}

// GetStats returns platform statistics
func (h *AdminHandler) GetStats(c *gin.Context) {
	db := database.GetDB()

	var stats struct {
		TotalUsers      int64 `json:"total_users"`
		TotalDevelopers int64 `json:"total_developers"`
		TotalCompanies  int64 `json:"total_companies"`
		TotalJobs       int64 `json:"total_jobs"`
		TotalReports    int64 `json:"total_reports"`
		PendingReviews  int64 `json:"pending_reviews"`
	}

	db.Model(&models.User{}).Count(&stats.TotalUsers)
	db.Model(&models.DeveloperProfile{}).Count(&stats.TotalDevelopers)
	db.Model(&models.CompanyProfile{}).Count(&stats.TotalCompanies)
	db.Model(&models.JobPost{}).Count(&stats.TotalJobs)
	db.Model(&models.UserReport{}).Count(&stats.TotalReports)
	db.Model(&models.CompanyReview{}).Where("is_approved = ?", false).Count(&stats.PendingReviews)

	c.JSON(http.StatusOK, gin.H{
		"message": "Statistics retrieved successfully",
		"data":    stats,
	})
}

// ListUsers returns all users with pagination
func (h *AdminHandler) ListUsers(c *gin.Context) {
	page, limit := getPaginationFromQuery(c)
	role := c.Query("role")

	users, total, err := h.userRepo.List(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	// Filter by role if specified
	if role != "" {
		var filtered []models.User
		for _, user := range users {
			if user.Role == role {
				filtered = append(filtered, user)
			}
		}
		users = filtered
		total = int64(len(filtered))
	}

	result := utils.PaginationResult{
		Data:       users,
		TotalCount: total,
		Page:       page,
		Limit:      limit,
		TotalPages: utils.CalculateTotalPages(total, limit),
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Users retrieved successfully",
		"data":    result,
	})
}

// ApproveReview approves a company review
func (h *AdminHandler) ApproveReview(c *gin.Context) {
	reviewID, err := getIDFromURL(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	review, err := h.companyRepo.FindReviewByID(reviewID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
		return
	}

	review.IsApproved = true
	if err := h.companyRepo.UpdateReview(review); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to approve review"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Review approved successfully",
		"data":    review,
	})
}

// RejectReview rejects/deletes a company review
func (h *AdminHandler) RejectReview(c *gin.Context) {
	reviewID, err := getIDFromURL(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	if err := h.companyRepo.DeleteReview(reviewID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reject review"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Review rejected successfully",
		"data":    nil,
	})
}

// ListPendingReviews returns all pending reviews
func (h *AdminHandler) ListPendingReviews(c *gin.Context) {
	page, limit := getPaginationFromQuery(c)

	db := database.GetDB()
	var reviews []models.CompanyReview
	var total int64

	offset := (page - 1) * limit

	db.Model(&models.CompanyReview{}).Where("is_approved = ?", false).Count(&total)
	err := db.Where("is_approved = ?", false).
		Offset(offset).Limit(limit).
		Preload("User").Preload("Company").
		Find(&reviews).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pending reviews"})
		return
	}

	result := utils.PaginationResult{
		Data:       reviews,
		TotalCount: total,
		Page:       page,
		Limit:      limit,
		TotalPages: utils.CalculateTotalPages(total, limit),
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Pending reviews retrieved successfully",
		"data":    result,
	})
}

// ApproveComment approves a review comment
func (h *AdminHandler) ApproveComment(c *gin.Context) {
	commentID, err := getIDFromURL(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	comment, err := h.companyRepo.FindReviewCommentByID(commentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	comment.IsApproved = true
	if err := h.companyRepo.UpdateReviewComment(comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to approve comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Comment approved successfully",
		"data":    comment,
	})
}

// ListReports returns all user reports
func (h *AdminHandler) ListReports(c *gin.Context) {
	page, limit := getPaginationFromQuery(c)
	status := c.Query("status")
	reportType := c.Query("type")

	reports, total, err := h.reportRepo.List(page, limit, status, reportType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reports"})
		return
	}

	result := utils.PaginationResult{
		Data:       reports,
		TotalCount: total,
		Page:       page,
		Limit:      limit,
		TotalPages: utils.CalculateTotalPages(total, limit),
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Reports retrieved successfully",
		"data":    result,
	})
}

// UpdateReportStatus updates the status of a user report
func (h *AdminHandler) UpdateReportStatus(c *gin.Context) {
	reportID, err := getIDFromURL(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid report ID"})
		return
	}

	adminID, _ := middleware.GetUserID(c)

	var req struct {
		Status string `json:"status" validate:"required,oneof=reviewed resolved dismissed"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	report, err := h.reportRepo.FindByID(reportID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Report not found"})
		return
	}

	report.Status = req.Status
	report.ReviewedBy = &adminID
	now := time.Now()
	report.ReviewedAt = &now

	if err := h.reportRepo.Update(report); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update report"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Report status updated successfully",
		"data":    report,
	})
}

// DeleteUser soft deletes a user (admin action)
func (h *AdminHandler) DeleteUser(c *gin.Context) {
	userID, err := getIDFromURL(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := h.userRepo.Delete(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
		"data":    nil,
	})
}
