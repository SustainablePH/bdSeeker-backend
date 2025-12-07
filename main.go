package main

import (
	"log"
	"net/http"

	"github.com/bishworup11/bdSeeker-backend/internal/config"
	"github.com/bishworup11/bdSeeker-backend/internal/database"
	"github.com/bishworup11/bdSeeker-backend/internal/handlers"
	"github.com/bishworup11/bdSeeker-backend/internal/middleware"
	"github.com/bishworup11/bdSeeker-backend/internal/repositories"
	"github.com/bishworup11/bdSeeker-backend/internal/services"
	"github.com/gorilla/mux"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	if err := database.Connect(cfg); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run migrations
	if err := database.Migrate(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Seed default admin user
	if err := database.SeedAdminUser(); err != nil {
		log.Fatalf("Failed to seed admin user: %v", err)
	}

	// Initialize repositories
	db := database.GetDB()
	userRepo := repositories.NewUserRepository(db)

	// Initialize services
	authService := services.NewAuthService(userRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	companyHandler := handlers.NewCompanyHandler()
	developerHandler := handlers.NewDeveloperHandler()
	jobHandler := handlers.NewJobHandler()
	techHandler := handlers.NewTechHandler()
	adminHandler := handlers.NewAdminHandler()

	// Setup router
	router := mux.NewRouter()

	// Apply global middleware
	router.Use(middleware.CORSMiddleware)
	router.Use(middleware.ErrorHandler)

	// API routes
	api := router.PathPrefix("/api/v1").Subrouter()

	// Health check
	api.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}).Methods("GET")

	// Auth routes (public)
	api.HandleFunc("/auth/register", authHandler.Register).Methods("POST")
	api.HandleFunc("/auth/login", authHandler.Login).Methods("POST")
	api.HandleFunc("/auth/logout", authHandler.Logout).Methods("POST")

	// Protected auth routes
	authRoutes := api.PathPrefix("/auth").Subrouter()
	authRoutes.Use(middleware.AuthMiddleware)
	authRoutes.HandleFunc("/me", authHandler.GetMe).Methods("GET")

	// Technology routes (public read, admin write)
	api.HandleFunc("/technologies", techHandler.ListTechnologies).Methods("GET")
	api.HandleFunc("/languages", techHandler.ListLanguages).Methods("GET")

	// Company routes
	api.HandleFunc("/companies", companyHandler.ListCompanies).Methods("GET")
	api.HandleFunc("/companies/{id}", companyHandler.GetCompany).Methods("GET")
	api.HandleFunc("/companies/{id}/reviews", companyHandler.ListReviews).Methods("GET")

	// Protected company routes
	companyRoutes := api.PathPrefix("/companies").Subrouter()
	companyRoutes.Use(middleware.AuthMiddleware)
	companyRoutes.HandleFunc("", companyHandler.CreateCompany).Methods("POST")
	companyRoutes.HandleFunc("/{id}/ratings", companyHandler.RateCompany).Methods("POST")
	companyRoutes.HandleFunc("/{id}/reviews", companyHandler.CreateReview).Methods("POST")

	// Developer routes
	api.HandleFunc("/developers", developerHandler.ListDevelopers).Methods("GET")
	api.HandleFunc("/developers/{id}", developerHandler.GetDeveloper).Methods("GET")

	// Protected developer routes
	devRoutes := api.PathPrefix("/developers").Subrouter()
	devRoutes.Use(middleware.AuthMiddleware)
	devRoutes.HandleFunc("", developerHandler.CreateDeveloper).Methods("POST")

	// Job routes
	api.HandleFunc("/jobs", jobHandler.ListJobs).Methods("GET")
	api.HandleFunc("/jobs/{id}", jobHandler.GetJob).Methods("GET")

	// Protected job routes
	jobRoutes := api.PathPrefix("/jobs").Subrouter()
	jobRoutes.Use(middleware.AuthMiddleware)
	jobRoutes.HandleFunc("", jobHandler.CreateJob).Methods("POST")
	jobRoutes.HandleFunc("/{id}/reactions", jobHandler.ReactToJob).Methods("POST")
	jobRoutes.HandleFunc("/{id}/comments", jobHandler.CommentOnJob).Methods("POST")

	// Admin routes (protected, admin only)
	adminRoutes := api.PathPrefix("/admin").Subrouter()
	adminRoutes.Use(middleware.AuthMiddleware)
	adminRoutes.Use(middleware.RoleMiddleware("admin"))
	
	// Admin - Statistics
	adminRoutes.HandleFunc("/stats", adminHandler.GetStats).Methods("GET")
	
	// Admin - User Management
	adminRoutes.HandleFunc("/users", adminHandler.ListUsers).Methods("GET")
	adminRoutes.HandleFunc("/users/{id}", adminHandler.DeleteUser).Methods("DELETE")
	
	// Admin - Review Management
	adminRoutes.HandleFunc("/reviews/pending", adminHandler.ListPendingReviews).Methods("GET")
	adminRoutes.HandleFunc("/reviews/{id}/approve", adminHandler.ApproveReview).Methods("PUT")
	adminRoutes.HandleFunc("/reviews/{id}/reject", adminHandler.RejectReview).Methods("DELETE")
	
	// Admin - Comment Management
	adminRoutes.HandleFunc("/comments/{id}/approve", adminHandler.ApproveComment).Methods("PUT")
	
	// Admin - Report Management
	adminRoutes.HandleFunc("/reports", adminHandler.ListReports).Methods("GET")
	adminRoutes.HandleFunc("/reports/{id}", adminHandler.UpdateReportStatus).Methods("PUT")

	// Start server
	addr := cfg.ServerHost + ":" + cfg.ServerPort
	log.Printf("🚀 Server starting on %s", addr)
	log.Printf("📚 API Documentation: http://%s/api/v1/health", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
