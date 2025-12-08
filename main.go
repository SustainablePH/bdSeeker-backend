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
	"github.com/gin-gonic/gin"
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

	// Setup Gin router
	// Use gin.New() for custom middleware control
	router := gin.New()

	// Apply global middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.ErrorHandler())

	// API routes
	api := router.Group("/api/v1")

	// Health check
	api.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Auth routes (public)
	api.POST("/auth/register", authHandler.Register)
	api.POST("/auth/login", authHandler.Login)
	api.POST("/auth/logout", authHandler.Logout)

	// Protected auth routes
	authRoutes := api.Group("/auth")
	authRoutes.Use(middleware.AuthMiddleware())
	{
		authRoutes.GET("/me", authHandler.GetMe)
	}

	// Technology routes (public read, admin write)
	api.GET("/technologies", techHandler.ListTechnologies)
	api.GET("/languages", techHandler.ListLanguages)

	// Company routes (public)
	api.GET("/companies", companyHandler.ListCompanies)
	api.GET("/companies/:id", companyHandler.GetCompany)
	api.GET("/companies/:id/reviews", companyHandler.ListReviews)

	// Protected company routes
	companyRoutes := api.Group("/companies")
	companyRoutes.Use(middleware.AuthMiddleware())
	{
		companyRoutes.POST("", companyHandler.CreateCompany)
		companyRoutes.POST("/:id/ratings", companyHandler.RateCompany)
		companyRoutes.POST("/:id/reviews", companyHandler.CreateReview)
	}

	// Developer routes (public)
	api.GET("/developers", developerHandler.ListDevelopers)
	api.GET("/developers/:id", developerHandler.GetDeveloper)

	// Protected developer routes
	devRoutes := api.Group("/developers")
	devRoutes.Use(middleware.AuthMiddleware())
	{
		devRoutes.POST("", developerHandler.CreateDeveloper)
	}

	// Job routes (public)
	api.GET("/jobs", jobHandler.ListJobs)
	api.GET("/jobs/:id", jobHandler.GetJob)

	// Protected job routes
	jobRoutes := api.Group("/jobs")
	jobRoutes.Use(middleware.AuthMiddleware())
	{
		jobRoutes.POST("", jobHandler.CreateJob)
		jobRoutes.POST("/:id/reactions", jobHandler.ReactToJob)
		jobRoutes.POST("/:id/comments", jobHandler.CommentOnJob)
	}

	// Admin routes (protected, admin only)
	adminRoutes := api.Group("/admin")
	adminRoutes.Use(middleware.AuthMiddleware())
	adminRoutes.Use(middleware.RoleMiddleware("admin"))
	{
		// Admin - Statistics
		adminRoutes.GET("/stats", adminHandler.GetStats)

		// Admin - User Management
		adminRoutes.GET("/users", adminHandler.ListUsers)
		adminRoutes.DELETE("/users/:id", adminHandler.DeleteUser)

		// Admin - Review Management
		adminRoutes.GET("/reviews/pending", adminHandler.ListPendingReviews)
		adminRoutes.PUT("/reviews/:id/approve", adminHandler.ApproveReview)
		adminRoutes.DELETE("/reviews/:id/reject", adminHandler.RejectReview)

		// Admin - Comment Management
		adminRoutes.PUT("/comments/:id/approve", adminHandler.ApproveComment)

		// Admin - Report Management
		adminRoutes.GET("/reports", adminHandler.ListReports)
		adminRoutes.PUT("/reports/:id", adminHandler.UpdateReportStatus)
	}

	// Start server
	addr := cfg.ServerHost + ":" + cfg.ServerPort
	log.Printf("🚀 Server starting on %s", addr)
	log.Printf("📚 API Documentation: http://%s/api/v1/health", addr)
	log.Printf("🔧 Environment: %s", cfg.Environment)
	
	if err := router.Run(addr); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
