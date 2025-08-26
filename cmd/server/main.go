package main

import (
	"go-starter/internal/modules/auth/models"
	"go-starter/internal/server"
	"go-starter/pkg/config"
	"go-starter/pkg/db"
	"log"
)

func main() {
	cfg := config.New()

	database, err := db.New(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		if err := database.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()

	// Run GORM auto-migration with cleanup
	if err := database.MigrateWithCleanup(&models.User{}); err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}
	log.Println("Database migrations completed successfully")

	srv := server.New(cfg, database)
	srv.Start()
}
