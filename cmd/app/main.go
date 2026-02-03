package main

import (
	"log"

	"freepass-2026/internal/handler/rest"
	"freepass-2026/internal/repository"
	"freepass-2026/internal/service"
	"freepass-2026/pkg/config"
	"freepass-2026/pkg/database/mariadb"
	"freepass-2026/pkg/jwt"
	"freepass-2026/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Setup Gin mode
	gin.SetMode(cfg.Server.Mode)

	// Initialize database connection
	db, err := mariadb.NewConnection(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto migrate database schemas
	if err := mariadb.AutoMigrate(db); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize JWT service
	jwtService := jwt.NewJWTService(cfg.JWT.Secret, cfg.JWT.ExpireHour)

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)

	// Initialize services
	userService := service.NewUserService(userRepo, jwtService)

	// Initialize handlers
	userHandler := rest.NewUserHandler(userService)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(jwtService)

	// Setup Gin router
	router := gin.Default()

	// API v1 routes
	api := router.Group("/api/v1")
	{
		userHandler.RegisterRoutes(api, authMiddleware)
	}

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "canteen-api",
		})
	})

	// Start server
	log.Printf("Server starting on port %s", cfg.Server.Port)
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
