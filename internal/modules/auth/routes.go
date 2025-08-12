package auth

import (
	"products/internal/modules/auth/handlers"
	"products/internal/modules/auth/middlewares"
	"products/internal/modules/auth/repositories"
	"products/internal/modules/auth/services"
	"products/pkg/db"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, api *echo.Group, database *db.Database) {

	authRepo := repositories.NewRepository(database)
	jwtService := services.NewJWTService()
	authService := services.NewAuthService(authRepo, jwtService)
	authHandler := handlers.NewHandler(authService)

	authGroup := api.Group("/auth")
	authGroup.POST("/register", authHandler.Register)
	authGroup.POST("/login", authHandler.Login)

	api.GET("/profile", authHandler.GetProfile, middlewares.JWTMiddleware(jwtService))
}
