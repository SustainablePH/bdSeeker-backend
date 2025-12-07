package middleware

import (
	"log"
	"net/http"

	"github.com/bishworup11/bdSeeker-backend/pkg/utils"
)

// ErrorHandler recovers from panics and handles errors
func ErrorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				utils.RespondError(w, http.StatusInternalServerError, "Internal server error")
			}
		}()

		next.ServeHTTP(w, r)
	})
}
