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
	authAPIHandler := handlers.NewAuthAPIHandler(authService)
	authWEBHandler := handlers.NewAuthWEBHandler(authService)

	e.GET("/login", authWEBHandler.ViewLogin)
	e.POST("/login", authWEBHandler.Login)
	e.GET("/register", authWEBHandler.ViewRegister)
	e.POST("/register", authWEBHandler.Register)

	authGroup := api.Group("/auth")
	authGroup.POST("/register", authAPIHandler.Register)
	authGroup.POST("/login", authAPIHandler.Login)

	api.GET("/profile", authAPIHandler.GetProfile, middlewares.JWTMiddleware(jwtService))
}
