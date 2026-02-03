package mariadb

import (
	"fmt"
	"log"

	"freepass-2026/entity"
	"freepass-2026/pkg/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewConnection(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName,
	)

	logLevel := logger.Info
	if cfg.Server.Mode == "release" {
		logLevel = logger.Error
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Database connection established")
	return db, nil
}

func AutoMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&entity.User{},
	)
	if err != nil {
		return fmt.Errorf("failed to auto migrate: %w", err)
	}

	log.Println("Database migration completed")
	return nil
}
