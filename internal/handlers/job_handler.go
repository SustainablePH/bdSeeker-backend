package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/bishworup11/bdSeeker-backend/internal/database"
	"github.com/bishworup11/bdSeeker-backend/internal/middleware"
	"github.com/bishworup11/bdSeeker-backend/internal/models"
	"github.com/bishworup11/bdSeeker-backend/internal/repositories"
	"github.com/bishworup11/bdSeeker-backend/pkg/utils"
)

type JobHandler struct {
	repo        *repositories.JobRepository
	companyRepo *repositories.CompanyRepository
}

func NewJobHandler() *JobHandler {
	db := database.GetDB()
	return &JobHandler{
		repo:        repositories.NewJobRepository(db),
		companyRepo: repositories.NewCompanyRepository(db),
	}
}

// ListJobs GET /api/v1/jobs
func (h *JobHandler) ListJobs(w http.ResponseWriter, r *http.Request) {
	page, limit := getPaginationFromQuery(r)
	
	filters := make(map[string]interface{})
	
	if workMode := r.URL.Query().Get("work_mode"); workMode != "" {
		filters["work_mode"] = workMode
	}
	if location := r.URL.Query().Get("location"); location != "" {
		filters["location"] = location
	}
	if search := r.URL.Query().Get("search"); search != "" {
		filters["search"] = search
	}
	if sortBy := r.URL.Query().Get("sort_by"); sortBy != "" {
		filters["sort_by"] = sortBy
	}

	jobs, total, err := h.repo.List(page, limit, filters)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to fetch jobs")
		return
	}

	result := utils.PaginationResult{
		Data:       jobs,
		TotalCount: total,
		Page:       page,
		Limit:      limit,
		TotalPages: utils.CalculateTotalPages(total, limit),
	}

	utils.RespondSuccess(w, "Jobs retrieved successfully", result)
}

// GetJob GET /api/v1/jobs/:id
func (h *JobHandler) GetJob(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromURL(r)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid job ID")
		return
	}

	job, err := h.repo.FindByID(id)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, "Job not found")
		return
	}

	utils.RespondSuccess(w, "Job retrieved successfully", job)
}

// CreateJob POST /api/v1/jobs
func (h *JobHandler) CreateJob(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserID(r)
	userRole, _ := middleware.GetUserRole(r)

	var req struct {
		Title              string  `json:"title" validate:"required"`
		Description        string  `json:"description" validate:"required"`
		SalaryMin          float64 `json:"salary_min"`
		SalaryMax          float64 `json:"salary_max"`
		ExperienceMinYears int     `json:"experience_min_years"`
		ExperienceMaxYears int     `json:"experience_max_years"`
		WorkMode           string  `json:"work_mode"`
		Location           string  `json:"location"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	var companyID uint
	if userRole != "admin" {
		// Get company profile for user
		company, err := h.companyRepo.FindByUserID(userID)
		if err != nil {
			utils.RespondError(w, http.StatusBadRequest, "User must have a company profile to post jobs")
			return
		}
		companyID = company.ID
	}

	job := &models.JobPost{
		CompanyID:          companyID,
		Title:              req.Title,
		Description:        req.Description,
		SalaryMin:          req.SalaryMin,
		SalaryMax:          req.SalaryMax,
		ExperienceMinYears: req.ExperienceMinYears,
		ExperienceMaxYears: req.ExperienceMaxYears,
		WorkMode:           req.WorkMode,
		Location:           req.Location,
	}

	if err := h.repo.Create(job); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to create job")
		return
	}

	utils.RespondCreated(w, "Job created successfully", job)
}

// ReactToJob POST /api/v1/jobs/:id/reactions
func (h *JobHandler) ReactToJob(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserID(r)
	jobID, err := getIDFromURL(r)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid job ID")
		return
	}

	var req struct {
		Type string `json:"type" validate:"required,oneof=like bookmark apply"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Check if reaction exists
	existing, err := h.repo.FindReaction(userID, jobID, "")
	if err == nil && existing != nil {
		existing.Type = req.Type
		h.repo.UpdateReaction(existing)
		utils.RespondSuccess(w, "Reaction updated successfully", existing)
		return
	}

	reaction := &models.PostReaction{
		UserID:    userID,
		JobPostID: jobID,
		Type:      req.Type,
	}

	if err := h.repo.CreateReaction(reaction); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to create reaction")
		return
	}

	utils.RespondCreated(w, "Reaction created successfully", reaction)
}

// CommentOnJob POST /api/v1/jobs/:id/comments
func (h *JobHandler) CommentOnJob(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserID(r)
	jobID, err := getIDFromURL(r)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid job ID")
		return
	}

	var req struct {
		Content string `json:"content" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	comment := &models.PostComment{
		UserID:    userID,
		JobPostID: jobID,
		Content:   req.Content,
	}

	if err := h.repo.CreateComment(comment); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to create comment")
		return
	}

	utils.RespondCreated(w, "Comment created successfully", comment)
}
