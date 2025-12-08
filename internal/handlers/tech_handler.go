package handlers

import (
	"net/http"

	"github.com/bishworup11/bdSeeker-backend/internal/database"
	"github.com/bishworup11/bdSeeker-backend/internal/repositories"
	"github.com/gin-gonic/gin"
)

type TechHandler struct {
	repo *repositories.TechRepository
}

func NewTechHandler() *TechHandler {
	return &TechHandler{
		repo: repositories.NewTechRepository(database.GetDB()),
	}
}

// ListTechnologies GET /api/v1/technologies
func (h *TechHandler) ListTechnologies(c *gin.Context) {
	search := c.Query("search")

	techs, err := h.repo.ListTechnologies(search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch technologies"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Technologies retrieved successfully",
		"data":    techs,
	})
}

// ListLanguages GET /api/v1/languages
func (h *TechHandler) ListLanguages(c *gin.Context) {
	search := c.Query("search")

	langs, err := h.repo.ListLanguages(search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch languages"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Languages retrieved successfully",
		"data":    langs,
	})
}
