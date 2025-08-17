package dashboard

import (
	"products/internal/modules/auth/middlewares"
	"products/internal/modules/auth/services"
	"products/internal/modules/dashboard/handlers"
	"products/pkg/db"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, api *echo.Group, database *db.Database) {
	dashboardHandler := handlers.NewDashboardWebHandler()
	jwtService := services.NewJWTService()

	e.GET("/dashboard", dashboardHandler.ViewDashboard, middlewares.JWTMiddleware(jwtService))
}
