package handlers

import (
	"net/http"

	"github.com/bishworup11/bdSeeker-backend/internal/database"
	"github.com/bishworup11/bdSeeker-backend/internal/middleware"
	"github.com/bishworup11/bdSeeker-backend/internal/models"
	"github.com/bishworup11/bdSeeker-backend/internal/repositories"
	"github.com/bishworup11/bdSeeker-backend/pkg/utils"
	"github.com/gin-gonic/gin"
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
func (h *JobHandler) ListJobs(c *gin.Context) {
	page, limit := getPaginationFromQuery(c)

	filters := make(map[string]interface{})

	if workMode := c.Query("work_mode"); workMode != "" {
		filters["work_mode"] = workMode
	}
	if location := c.Query("location"); location != "" {
		filters["location"] = location
	}
	if search := c.Query("search"); search != "" {
		filters["search"] = search
	}
	if sortBy := c.Query("sort_by"); sortBy != "" {
		filters["sort_by"] = sortBy
	}

	jobs, total, err := h.repo.List(page, limit, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch jobs"})
		return
	}

	result := utils.PaginationResult{
		Data:       jobs,
		TotalCount: total,
		Page:       page,
		Limit:      limit,
		TotalPages: utils.CalculateTotalPages(total, limit),
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Jobs retrieved successfully",
		"data":    result,
	})
}

// GetJob GET /api/v1/jobs/:id
func (h *JobHandler) GetJob(c *gin.Context) {
	id, err := getIDFromURL(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job ID"})
		return
	}

	job, err := h.repo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Job retrieved successfully",
		"data":    job,
	})
}

// CreateJob POST /api/v1/jobs
func (h *JobHandler) CreateJob(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	userRole, _ := middleware.GetUserRole(c)

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

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	var companyID uint
	if userRole != "admin" {
		// Get company profile for user
		company, err := h.companyRepo.FindByUserID(userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User must have a company profile to post jobs"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create job"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Job created successfully",
		"data":    job,
	})
}

// ReactToJob POST /api/v1/jobs/:id/reactions
func (h *JobHandler) ReactToJob(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	jobID, err := getIDFromURL(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job ID"})
		return
	}

	var req struct {
		Type string `json:"type" validate:"required,oneof=like bookmark apply"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Check if reaction exists
	existing, err := h.repo.FindReaction(userID, jobID, "")
	if err == nil && existing != nil {
		existing.Type = req.Type
		h.repo.UpdateReaction(existing)
		c.JSON(http.StatusOK, gin.H{
			"message": "Reaction updated successfully",
			"data":    existing,
		})
		return
	}

	reaction := &models.PostReaction{
		UserID:    userID,
		JobPostID: jobID,
		Type:      req.Type,
	}

	if err := h.repo.CreateReaction(reaction); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create reaction"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Reaction created successfully",
		"data":    reaction,
	})
}

// CommentOnJob POST /api/v1/jobs/:id/comments
func (h *JobHandler) CommentOnJob(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	jobID, err := getIDFromURL(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job ID"})
		return
	}

	var req struct {
		Content string `json:"content" validate:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	comment := &models.PostComment{
		UserID:    userID,
		JobPostID: jobID,
		Content:   req.Content,
	}

	if err := h.repo.CreateComment(comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Comment created successfully",
		"data":    comment,
	})
}
