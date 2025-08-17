package main

import (
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
	defer database.Close()

	srv := server.New(cfg, database)
	srv.Start()
}
