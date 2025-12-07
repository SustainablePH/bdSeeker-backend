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

type DeveloperHandler struct {
	repo *repositories.DeveloperRepository
}

func NewDeveloperHandler() *DeveloperHandler {
	return &DeveloperHandler{
		repo: repositories.NewDeveloperRepository(database.GetDB()),
	}
}

// ListDevelopers GET /api/v1/developers
func (h *DeveloperHandler) ListDevelopers(w http.ResponseWriter, r *http.Request) {
	page, limit := getPaginationFromQuery(r)
	
	var techIDs, langIDs []uint
	
	developers, total, err := h.repo.List(page, limit, techIDs, langIDs)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to fetch developers")
		return
	}

	result := utils.PaginationResult{
		Data:       developers,
		TotalCount: total,
		Page:       page,
		Limit:      limit,
		TotalPages: utils.CalculateTotalPages(total, limit),
	}

	utils.RespondSuccess(w, "Developers retrieved successfully", result)
}

// GetDeveloper GET /api/v1/developers/:id
func (h *DeveloperHandler) GetDeveloper(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromURL(r)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid developer ID")
		return
	}

	developer, err := h.repo.FindByID(id)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, "Developer not found")
		return
	}

	utils.RespondSuccess(w, "Developer retrieved successfully", developer)
}

// CreateDeveloper POST /api/v1/developers
func (h *DeveloperHandler) CreateDeveloper(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserID(r)

	var req struct {
		Bio string `json:"bio"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Check if user already has a developer profile
	existing, err := h.repo.FindByUserID(userID)
	if err == nil && existing.ID > 0 {
		utils.RespondError(w, http.StatusBadRequest, "User already has a developer profile")
		return
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to check existing profile")
		return
	}

	developer := &models.DeveloperProfile{
		UserID: userID,
		Bio:    req.Bio,
	}

	if err := h.repo.Create(developer); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to create developer profile")
		return
	}

	utils.RespondCreated(w, "Developer profile created successfully", developer)
}
