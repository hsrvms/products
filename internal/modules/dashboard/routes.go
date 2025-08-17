package dashboard

import (
	"go-starter/internal/modules/auth/middlewares"
	"go-starter/internal/modules/auth/services"
	"go-starter/internal/modules/dashboard/handlers"
	"go-starter/pkg/db"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, api *echo.Group, database *db.Database) {
	dashboardHandler := handlers.NewDashboardWebHandler()
	jwtService := services.NewJWTService()

	e.GET("/dashboard", dashboardHandler.ViewDashboard, middlewares.JWTMiddleware(jwtService))
}
