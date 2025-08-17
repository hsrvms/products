package handlers

import (
	"go-starter/internal/modules/dashboard/views"

	"github.com/labstack/echo/v4"
)

type DashboardWEBHandler struct{}

func NewDashboardWebHandler() *DashboardWEBHandler {
	return &DashboardWEBHandler{}
}

func (h *DashboardWEBHandler) ViewDashboard(c echo.Context) error {
	component := views.HomePage()
	return component.Render(c.Request().Context(), c.Response().Writer)
}
