package server

import (
	"context"
	"fmt"
	"go-starter/pkg/config"
	"go-starter/pkg/db"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	prettylogger "github.com/rdbell/echo-pretty-logger"
)

type Server struct {
	Echo   *echo.Echo
	DB     *db.Database
	Config *config.Config
}

func New(cfg *config.Config, database *db.Database) *Server {
	e := echo.New()

	e.Use(prettylogger.Logger)
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.Static("/assets", "web/assets")

	server := &Server{
		Echo:   e,
		DB:     database,
		Config: cfg,
	}

	server.initRoutes()

	return server
}

func (s *Server) Start() {
	addr := fmt.Sprintf(":%d", s.Config.Server.Port)

	httpServer := &http.Server{
		Addr:         addr,
		ReadTimeout:  s.Config.Server.ReadTimeout,
		WriteTimeout: s.Config.Server.WriteTimeout,
		IdleTimeout:  s.Config.Server.IdleTimeout,
	}

	go func() {
		if err := s.Echo.StartServer(httpServer); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Printf("Server started on %s", addr)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.Echo.Shutdown(ctx); err != nil {
		log.Fatalf("Failed to gracefully shutdown server: %v", err)
	}

	log.Println("Server stopped")
}
