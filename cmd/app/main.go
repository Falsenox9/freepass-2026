package main

import (
    "log"

    "freepass-2026/internal/handler/rest"
    "freepass-2026/internal/repository"
    "freepass-2026/internal/service"
    "freepass-2026/pkg/bcrypt"
    "freepass-2026/pkg/config"
    "freepass-2026/pkg/database/mariadb"
    "freepass-2026/pkg/jwt"
    "freepass-2026/pkg/middleware"
)

func main() {
    // Load config
    cfg := config.Load()

    // Initialize database
    db, err := mariadb.NewConnection(cfg)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

    // Auto migrate
    if err := mariadb.AutoMigrate(db); err != nil {
        log.Fatalf("Failed to migrate database: %v", err)
    }

    // Initialize bcrypt
    bcryptService := bcrypt.Init()

    // Initialize JWT
    jwtService := jwt.Init()

    // Initialize repositories
    repo := repository.NewRepository(db)

    // Initialize services
    svc := service.NewService(repo, bcryptService, jwtService)

    // Initialize middleware
    mw := middleware.Init(svc, jwtService)

    // Initialize REST handler
    restHandler := rest.NewRest(svc, mw)

    // Mount endpoints
    restHandler.MountEndPoint()

    // Run server
    restHandler.Run()
}