package dashboard

import (
	"products/internal/modules/dashboard/handlers"
	"products/pkg/db"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, api *echo.Group, database *db.Database) {
	dashboardHandler := handlers.NewHandler()

	e.GET("/", dashboardHandler.ViewDashboard)
}
