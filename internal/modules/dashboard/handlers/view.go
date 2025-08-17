package handlers

import (
	"products/web/templates/dashboard"

	"github.com/labstack/echo/v4"
)

func (h *DashboardHandler) ViewDashboard(c echo.Context) error {
	component := dashboard.HomePage()
	return component.Render(c.Request().Context(), c.Response().Writer)
}
