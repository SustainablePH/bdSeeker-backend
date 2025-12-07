package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/bishworup11/bdSeeker-backend/internal/database"
	"github.com/bishworup11/bdSeeker-backend/internal/middleware"
	"github.com/bishworup11/bdSeeker-backend/internal/models"
	"github.com/bishworup11/bdSeeker-backend/internal/repositories"
	"github.com/bishworup11/bdSeeker-backend/pkg/utils"
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
func (h *CompanyHandler) ListCompanies(w http.ResponseWriter, r *http.Request) {
	page, limit := getPaginationFromQuery(r)
	location := r.URL.Query().Get("location")
	
	var techIDs []uint
	if techIDsStr := r.URL.Query().Get("tech_ids"); techIDsStr != "" {
		// Parse comma-separated tech IDs
	}

	companies, total, err := h.repo.List(page, limit, location, techIDs)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to fetch companies")
		return
	}

	result := utils.PaginationResult{
		Data:       companies,
		TotalCount: total,
		Page:       page,
		Limit:      limit,
		TotalPages: utils.CalculateTotalPages(total, limit),
	}

	utils.RespondSuccess(w, "Companies retrieved successfully", result)
}

// GetCompany GET /api/v1/companies/:id
func (h *CompanyHandler) GetCompany(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromURL(r)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid company ID")
		return
	}

	company, err := h.repo.FindByID(id)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, "Company not found")
		return
	}

	utils.RespondSuccess(w, "Company retrieved successfully", company)
}

// CreateCompany POST /api/v1/companies
func (h *CompanyHandler) CreateCompany(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserID(r)

	var req struct {
		CompanyName   string `json:"company_name" validate:"required"`
		Description   string `json:"description"`
		Website       string `json:"website"`
		Location      string `json:"location"`
		TechnologyIDs []uint `json:"technology_ids"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Check if user already has a company profile
	existing, err := h.repo.FindByUserID(userID)
	if err == nil && existing.ID > 0 {
		utils.RespondError(w, http.StatusBadRequest, "User already has a company profile")
		return
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to check existing profile")
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
		utils.RespondError(w, http.StatusInternalServerError, "Failed to create company")
		return
	}

	utils.RespondCreated(w, "Company created successfully", company)
}

// RateCompany POST /api/v1/companies/:id/ratings
func (h *CompanyHandler) RateCompany(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserID(r)
	companyID, err := getIDFromURL(r)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid company ID")
		return
	}

	var req struct {
		Rating int `json:"rating" validate:"required,min=1,max=5"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Check if rating exists
	existing, err := h.repo.FindRatingByUserAndCompany(userID, companyID)
	if err == nil && existing != nil {
		// Update existing rating
		existing.Rating = req.Rating
		h.repo.UpdateRating(existing)
		utils.RespondSuccess(w, "Rating updated successfully", existing)
		return
	}

	// Create new rating
	rating := &models.CompanyRating{
		CompanyID: companyID,
		UserID:    userID,
		Rating:    req.Rating,
	}

	if err := h.repo.CreateRating(rating); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to create rating")
		return
	}

	utils.RespondCreated(w, "Rating created successfully", rating)
}

// CreateReview POST /api/v1/companies/:id/reviews
func (h *CompanyHandler) CreateReview(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserID(r)
	companyID, err := getIDFromURL(r)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid company ID")
		return
	}

	var req struct {
		Content string `json:"content" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	review := &models.CompanyReview{
		CompanyID:  companyID,
		UserID:     userID,
		Content:    req.Content,
		IsApproved: false,
	}

	if err := h.repo.CreateReview(review); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to create review")
		return
	}

	utils.RespondCreated(w, "Review created successfully (pending approval)", review)
}

// ListReviews GET /api/v1/companies/:id/reviews
func (h *CompanyHandler) ListReviews(w http.ResponseWriter, r *http.Request) {
	companyID, err := getIDFromURL(r)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid company ID")
		return
	}

	page, limit := getPaginationFromQuery(r)
	reviews, total, err := h.repo.ListReviews(companyID, page, limit, true)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to fetch reviews")
		return
	}

	result := utils.PaginationResult{
		Data:       reviews,
		TotalCount: total,
		Page:       page,
		Limit:      limit,
		TotalPages: utils.CalculateTotalPages(total, limit),
	}

	utils.RespondSuccess(w, "Reviews retrieved successfully", result)
}
