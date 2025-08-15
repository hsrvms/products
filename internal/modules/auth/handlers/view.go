package handlers

import (
	"products/web/templates/auth"

	"github.com/labstack/echo/v4"
)

func (h *AuthHandler) ViewLogin(c echo.Context) error {
	component := auth.LoginPage()
	return component.Render(c.Request().Context(), c.Response().Writer)
}

func (h *AuthHandler) ViewRegister(c echo.Context) error {
	component := auth.RegisterPage()
	return component.Render(c.Request().Context(), c.Response().Writer)
}
