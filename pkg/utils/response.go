package utils

import (
	"encoding/json"
	"net/http"
)

// APIResponse represents a standardized API response
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// PaginationParams represents pagination parameters
type PaginationParams struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

// PaginationResult represents paginated data with metadata
type PaginationResult struct {
	Data       interface{} `json:"data"`
	TotalCount int64       `json:"total_count"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalPages int         `json:"total_pages"`
}

// RespondJSON sends a JSON response
func RespondJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// RespondSuccess sends a success response
func RespondSuccess(w http.ResponseWriter, message string, data interface{}) {
	RespondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// RespondCreated sends a created response
func RespondCreated(w http.ResponseWriter, message string, data interface{}) {
	RespondJSON(w, http.StatusCreated, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// RespondError sends an error response
func RespondError(w http.ResponseWriter, statusCode int, message string) {
	RespondJSON(w, statusCode, APIResponse{
		Success: false,
		Error:   message,
	})
}

// RespondValidationError sends a validation error response
func RespondValidationError(w http.ResponseWriter, errors map[string]string) {
	RespondJSON(w, http.StatusBadRequest, APIResponse{
		Success: false,
		Error:   "Validation failed",
		Data:    errors,
	})
}

// CalculateTotalPages calculates the total number of pages
func CalculateTotalPages(totalCount int64, limit int) int {
	if limit == 0 {
		return 0
	}
	pages := int(totalCount) / limit
	if int(totalCount)%limit > 0 {
		pages++
	}
	return pages
}

// GetPaginationParams extracts pagination parameters from query string with defaults
func GetPaginationParams(page, limit int) PaginationParams {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	return PaginationParams{
		Page:  page,
		Limit: limit,
	}
}
