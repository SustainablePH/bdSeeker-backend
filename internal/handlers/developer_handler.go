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

type DeveloperHandler struct {
	repo *repositories.DeveloperRepository
}

func NewDeveloperHandler() *DeveloperHandler {
	return &DeveloperHandler{
		repo: repositories.NewDeveloperRepository(database.GetDB()),
	}
}

// ListDevelopers GET /api/v1/developers
func (h *DeveloperHandler) ListDevelopers(c *gin.Context) {
	page, limit := getPaginationFromQuery(c)

	var techIDs, langIDs []uint

	developers, total, err := h.repo.List(page, limit, techIDs, langIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch developers"})
		return
	}

	result := utils.PaginationResult{
		Data:       developers,
		TotalCount: total,
		Page:       page,
		Limit:      limit,
		TotalPages: utils.CalculateTotalPages(total, limit),
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Developers retrieved successfully",
		"data":    result,
	})
}

// GetDeveloper GET /api/v1/developers/:id
func (h *DeveloperHandler) GetDeveloper(c *gin.Context) {
	id, err := getIDFromURL(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid developer ID"})
		return
	}

	developer, err := h.repo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Developer not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Developer retrieved successfully",
		"data":    developer,
	})
}

// CreateDeveloper POST /api/v1/developers
func (h *DeveloperHandler) CreateDeveloper(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	var req struct {
		Bio string `json:"bio"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Check if user already has a developer profile
	existing, err := h.repo.FindByUserID(userID)
	if err == nil && existing.ID > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already has a developer profile"})
		return
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing profile"})
		return
	}

	developer := &models.DeveloperProfile{
		UserID: userID,
		Bio:    req.Bio,
	}

	if err := h.repo.Create(developer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create developer profile"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Developer profile created successfully",
		"data":    developer,
	})
}
