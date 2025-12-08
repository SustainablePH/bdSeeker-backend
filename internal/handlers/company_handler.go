package handlers

import (
	"errors"
	"net/http"

	"github.com/bishworup11/bdSeeker-backend/internal/database"
	"github.com/bishworup11/bdSeeker-backend/internal/middleware"
	"github.com/bishworup11/bdSeeker-backend/internal/models"
	"github.com/bishworup11/bdSeeker-backend/internal/repositories"
	"github.com/bishworup11/bdSeeker-backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CompanyHandler struct {
	repo *repositories.CompanyRepository
}

func NewCompanyHandler() *CompanyHandler {
	return &CompanyHandler{
		repo: repositories.NewCompanyRepository(database.GetDB()),
	}
}

// ListCompanies GET /api/v1/companies
func (h *CompanyHandler) ListCompanies(c *gin.Context) {
	page, limit := getPaginationFromQuery(c)
	location := c.Query("location")

	var techIDs []uint
	// Parse comma-separated tech IDs if needed
	
	companies, total, err := h.repo.List(page, limit, location, techIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch companies"})
		return
	}

	result := utils.PaginationResult{
		Data:       companies,
		TotalCount: total,
		Page:       page,
		Limit:      limit,
		TotalPages: utils.CalculateTotalPages(total, limit),
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Companies retrieved successfully",
		"data":    result,
	})
}

// GetCompany GET /api/v1/companies/:id
func (h *CompanyHandler) GetCompany(c *gin.Context) {
	id, err := getIDFromURL(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID"})
		return
	}

	company, err := h.repo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Company retrieved successfully",
		"data":    company,
	})
}

// CreateCompany POST /api/v1/companies
func (h *CompanyHandler) CreateCompany(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	var req struct {
		CompanyName   string `json:"company_name" validate:"required"`
		Description   string `json:"description"`
		Website       string `json:"website"`
		Location      string `json:"location"`
		TechnologyIDs []uint `json:"technology_ids"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Check if user already has a company profile
	existing, err := h.repo.FindByUserID(userID)
	if err == nil && existing.ID > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already has a company profile"})
		return
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing profile"})
		return
	}

	company := &models.CompanyProfile{
		UserID:      userID,
		CompanyName: req.CompanyName,
		Description: req.Description,
		Website:     req.Website,
		Location:    req.Location,
	}

	if err := h.repo.Create(company); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create company"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Company created successfully",
		"data":    company,
	})
}

// RateCompany POST /api/v1/companies/:id/ratings
func (h *CompanyHandler) RateCompany(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	companyID, err := getIDFromURL(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID"})
		return
	}

	var req struct {
		Rating int `json:"rating" validate:"required,min=1,max=5"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Check if rating exists
	existing, err := h.repo.FindRatingByUserAndCompany(userID, companyID)
	if err == nil && existing != nil {
		// Update existing rating
		existing.Rating = req.Rating
		h.repo.UpdateRating(existing)
		c.JSON(http.StatusOK, gin.H{
			"message": "Rating updated successfully",
			"data":    existing,
		})
		return
	}

	// Create new rating
	rating := &models.CompanyRating{
		CompanyID: companyID,
		UserID:    userID,
		Rating:    req.Rating,
	}

	if err := h.repo.CreateRating(rating); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create rating"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Rating created successfully",
		"data":    rating,
	})
}

// CreateReview POST /api/v1/companies/:id/reviews
func (h *CompanyHandler) CreateReview(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	companyID, err := getIDFromURL(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID"})
		return
	}

	var req struct {
		Content string `json:"content" validate:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	review := &models.CompanyReview{
		CompanyID:  companyID,
		UserID:     userID,
		Content:    req.Content,
		IsApproved: false,
	}

	if err := h.repo.CreateReview(review); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create review"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Review created successfully (pending approval)",
		"data":    review,
	})
}

// ListReviews GET /api/v1/companies/:id/reviews
func (h *CompanyHandler) ListReviews(c *gin.Context) {
	companyID, err := getIDFromURL(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID"})
		return
	}

	page, limit := getPaginationFromQuery(c)
	reviews, total, err := h.repo.ListReviews(companyID, page, limit, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reviews"})
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
		"message": "Reviews retrieved successfully",
		"data":    result,
	})
}
