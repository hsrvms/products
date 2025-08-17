package handlers

import (
	"net/http"
	"products/internal/modules/auth/dto"
	"products/internal/modules/auth/services"
	"products/internal/modules/auth/views"
	"products/internal/modules/auth/views/errors"

	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type AuthWEBHandler struct {
	authService *services.AuthService
	validator   *validator.Validate
}

func NewAuthWEBHandler(authService *services.AuthService) *AuthWEBHandler {
	return &AuthWEBHandler{
		authService: authService,
		validator:   validator.New(),
	}
}

func (h *AuthWEBHandler) ViewLogin(c echo.Context) error {
	component := views.LoginPage()
	return component.Render(c.Request().Context(), c.Response().Writer)
}

func (h *AuthWEBHandler) ViewRegister(c echo.Context) error {
	component := views.RegisterPage()
	return component.Render(c.Request().Context(), c.Response().Writer)
}

func (h *AuthWEBHandler) Login(c echo.Context) error {
	var req dto.LoginRequest

	req.Email = strings.ToLower(strings.TrimSpace(c.FormValue("email")))
	req.Password = c.FormValue("password")

	if err := h.validator.Struct(req); err != nil {
		c.Response().Header().Set("Content-Type", "text/html")
		c.Response().WriteHeader(http.StatusBadRequest)
		return errors.LoginError("Invalid email or password").Render(c.Request().Context(), c.Response().Writer)
	}

	response, err := h.authService.Login(c.Request().Context(), &req)
	if err != nil {

		if strings.Contains(err.Error(), "invalid credentials") || strings.Contains(err.Error(), "user not found") {
			c.Response().Header().Set("Content-Type", "text/html")
			c.Response().WriteHeader(http.StatusUnauthorized)
			return errors.LoginError("Invalid email or password").Render(c.Request().Context(), c.Response().Writer)
		}
		c.Response().Header().Set("Content-Type", "text/html")
		c.Response().WriteHeader(http.StatusInternalServerError)
		return errors.LoginError("Failed to login").Render(c.Request().Context(), c.Response().Writer)
	}

	cookie := new(http.Cookie)
	cookie.Name = "auth_token"
	cookie.Value = response.Token
	cookie.Path = "/"
	cookie.HttpOnly = true
	cookie.Secure = false // make it true for HTTPS
	cookie.SameSite = http.SameSiteStrictMode
	cookie.MaxAge = 3600

	c.SetCookie(cookie)
	c.Response().Header().Set("HX-Location", "/")

	return c.NoContent(http.StatusOK)

}

func (h *AuthWEBHandler) Register(c echo.Context) error {
	var req dto.RegisterRequest

	req.FirstName = strings.ToLower(strings.TrimSpace(c.FormValue("firstname")))
	req.LastName = strings.ToLower(strings.TrimSpace(c.FormValue("lastname")))
	req.Email = strings.ToLower(strings.TrimSpace(c.FormValue("email")))
	req.Password = c.FormValue("password")

	if err := h.validator.Struct(req); err != nil {
		c.Response().Header().Set("Content-Type", "text/html")
		c.Response().WriteHeader(http.StatusUnauthorized)
		return errors.LoginError("Invalid email or password").Render(c.Request().Context(), c.Response().Writer)
	}

	response, err := h.authService.Register(c.Request().Context(), &req)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			c.Response().Header().Set("Content-Type", "text/html")
			c.Response().WriteHeader(http.StatusUnauthorized)
			return errors.LoginError("User with this email already exists").Render(c.Request().Context(), c.Response().Writer)
		}
		c.Response().Header().Set("Content-Type", "text/html")
		c.Response().WriteHeader(http.StatusInternalServerError)
		return errors.LoginError("Failed to register").Render(c.Request().Context(), c.Response().Writer)
	}

	cookie := new(http.Cookie)
	cookie.Name = "auth_token"
	cookie.Value = response.Token
	cookie.Path = "/"
	cookie.HttpOnly = true
	cookie.Secure = false // make it true on https
	cookie.SameSite = http.SameSiteStrictMode
	cookie.MaxAge = 3600

	c.SetCookie(cookie)
	c.Response().Header().Set("HX-Location", "/")

	return c.NoContent(http.StatusCreated)
}
