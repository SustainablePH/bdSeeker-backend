package handlers

import (
	"net/http"

	"github.com/bishworup11/bdSeeker-backend/internal/database"
	"github.com/bishworup11/bdSeeker-backend/internal/repositories"
	"github.com/bishworup11/bdSeeker-backend/pkg/utils"
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
func (h *TechHandler) ListTechnologies(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	
	techs, err := h.repo.ListTechnologies(search)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to fetch technologies")
		return
	}

	utils.RespondSuccess(w, "Technologies retrieved successfully", techs)
}

// ListLanguages GET /api/v1/languages
func (h *TechHandler) ListLanguages(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	
	langs, err := h.repo.ListLanguages(search)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to fetch languages")
		return
	}

	utils.RespondSuccess(w, "Languages retrieved successfully", langs)
}
