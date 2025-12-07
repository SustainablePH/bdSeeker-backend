package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/bishworup11/bdSeeker-backend/internal/middleware"
	"github.com/bishworup11/bdSeeker-backend/internal/services"
	"github.com/bishworup11/bdSeeker-backend/pkg/utils"
	"github.com/gorilla/mux"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register handles user registration
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req services.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if errors := utils.ValidateStruct(&req); errors != nil {
		utils.RespondValidationError(w, errors)
		return
	}

	// Register user
	response, err := h.authService.Register(&req)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Set HTTP-only cookie for browser clients
	utils.SetAuthCookie(w, response.Token, 24*60*60) // 24 hours

	utils.RespondCreated(w, "User registered successfully", response)
}

// Login handles user authentication
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req services.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if errors := utils.ValidateStruct(&req); errors != nil {
		utils.RespondValidationError(w, errors)
		return
	}

	// Login user
	response, err := h.authService.Login(&req)
	if err != nil {
		utils.RespondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// Set HTTP-only cookie for browser clients
	utils.SetAuthCookie(w, response.Token, 24*60*60) // 24 hours

	utils.RespondSuccess(w, "Login successful", response)
}

// GetMe returns the current authenticated user's information
func (h *AuthHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		utils.RespondError(w, http.StatusUnauthorized, "User ID not found in context")
		return
	}

	user, err := h.authService.GetUserByID(userID)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, "User not found")
		return
	}

	utils.RespondSuccess(w, "User retrieved successfully", user)
}

// Logout handles user logout by clearing cookies
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Clear authentication cookie
	utils.ClearAuthCookie(w)
	utils.ClearRefreshCookie(w)
	utils.ClearSessionCookie(w)

	utils.RespondSuccess(w, "Logged out successfully", nil)
}

// Helper function to get pagination params from query string
func getPaginationFromQuery(r *http.Request) (int, int) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	
	params := utils.GetPaginationParams(page, limit)
	return params.Page, params.Limit
}

// Helper function to get ID from URL params
func getIDFromURL(r *http.Request) (uint, error) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}
