package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/bishworup11/bdSeeker-backend/internal/database"
	"github.com/bishworup11/bdSeeker-backend/internal/middleware"
	"github.com/bishworup11/bdSeeker-backend/internal/models"
	"github.com/bishworup11/bdSeeker-backend/internal/repositories"
	"github.com/bishworup11/bdSeeker-backend/pkg/utils"
)

type AdminHandler struct {
	userRepo     *repositories.UserRepository
	companyRepo  *repositories.CompanyRepository
	reportRepo   *repositories.ReportRepository
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
func (h *AdminHandler) GetStats(w http.ResponseWriter, r *http.Request) {
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

	utils.RespondSuccess(w, "Statistics retrieved successfully", stats)
}

// ListUsers returns all users with pagination
func (h *AdminHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	page, limit := getPaginationFromQuery(r)
	role := r.URL.Query().Get("role")

	users, total, err := h.userRepo.List(page, limit)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to fetch users")
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

	utils.RespondSuccess(w, "Users retrieved successfully", result)
}

// ApproveReview approves a company review
func (h *AdminHandler) ApproveReview(w http.ResponseWriter, r *http.Request) {
	reviewID, err := getIDFromURL(r)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid review ID")
		return
	}

	review, err := h.companyRepo.FindReviewByID(reviewID)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, "Review not found")
		return
	}

	review.IsApproved = true
	if err := h.companyRepo.UpdateReview(review); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to approve review")
		return
	}

	utils.RespondSuccess(w, "Review approved successfully", review)
}

// RejectReview rejects/deletes a company review
func (h *AdminHandler) RejectReview(w http.ResponseWriter, r *http.Request) {
	reviewID, err := getIDFromURL(r)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid review ID")
		return
	}

	if err := h.companyRepo.DeleteReview(reviewID); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to reject review")
		return
	}

	utils.RespondSuccess(w, "Review rejected successfully", nil)
}

// ListPendingReviews returns all pending reviews
func (h *AdminHandler) ListPendingReviews(w http.ResponseWriter, r *http.Request) {
	page, limit := getPaginationFromQuery(r)
	
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
		utils.RespondError(w, http.StatusInternalServerError, "Failed to fetch pending reviews")
		return
	}

	result := utils.PaginationResult{
		Data:       reviews,
		TotalCount: total,
		Page:       page,
		Limit:      limit,
		TotalPages: utils.CalculateTotalPages(total, limit),
	}

	utils.RespondSuccess(w, "Pending reviews retrieved successfully", result)
}

// ApproveComment approves a review comment
func (h *AdminHandler) ApproveComment(w http.ResponseWriter, r *http.Request) {
	commentID, err := getIDFromURL(r)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid comment ID")
		return
	}

	comment, err := h.companyRepo.FindReviewCommentByID(commentID)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, "Comment not found")
		return
	}

	comment.IsApproved = true
	if err := h.companyRepo.UpdateReviewComment(comment); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to approve comment")
		return
	}

	utils.RespondSuccess(w, "Comment approved successfully", comment)
}

// ListReports returns all user reports
func (h *AdminHandler) ListReports(w http.ResponseWriter, r *http.Request) {
	page, limit := getPaginationFromQuery(r)
	status := r.URL.Query().Get("status")
	reportType := r.URL.Query().Get("type")

	reports, total, err := h.reportRepo.List(page, limit, status, reportType)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to fetch reports")
		return
	}

	result := utils.PaginationResult{
		Data:       reports,
		TotalCount: total,
		Page:       page,
		Limit:      limit,
		TotalPages: utils.CalculateTotalPages(total, limit),
	}

	utils.RespondSuccess(w, "Reports retrieved successfully", result)
}

// UpdateReportStatus updates the status of a user report
func (h *AdminHandler) UpdateReportStatus(w http.ResponseWriter, r *http.Request) {
	reportID, err := getIDFromURL(r)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid report ID")
		return
	}

	adminID, _ := middleware.GetUserID(r)

	var req struct {
		Status string `json:"status" validate:"required,oneof=reviewed resolved dismissed"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	report, err := h.reportRepo.FindByID(reportID)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, "Report not found")
		return
	}

	report.Status = req.Status
	report.ReviewedBy = &adminID
	now := time.Now()
	report.ReviewedAt = &now

	if err := h.reportRepo.Update(report); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to update report")
		return
	}

	utils.RespondSuccess(w, "Report status updated successfully", report)
}

// DeleteUser soft deletes a user (admin action)
func (h *AdminHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID, err := getIDFromURL(r)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	if err := h.userRepo.Delete(userID); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to delete user")
		return
	}

	utils.RespondSuccess(w, "User deleted successfully", nil)
}
