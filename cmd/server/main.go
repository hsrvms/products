package main

import (
	"log"
	"products/internal/server"
	"products/pkg/config"
	"products/pkg/db"
)

func main() {
	cfg := config.New()

	database, err := db.New(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	srv := server.New(cfg, database)
	srv.Start()
}
