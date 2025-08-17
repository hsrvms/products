package handlers

import (
	"go-starter/internal/modules/auth/dto"
	"go-starter/internal/modules/auth/services"
	"go-starter/internal/modules/auth/views"
	"go-starter/internal/modules/auth/views/errors"
	"net/http"

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

	req := dto.LoginRequest{
		Email:    strings.ToLower(strings.TrimSpace(c.FormValue("email"))),
		Password: c.FormValue("password"),
	}

	if err := h.validator.Struct(req); err != nil {
		return h.renderError(c, http.StatusBadRequest, "Invalid email or password")
	}

	response, err := h.authService.Login(c.Request().Context(), &req)
	if err != nil {
		statusCode, message := h.categorizeAuthError(err)
		return h.renderError(c, statusCode, message)
	}

	return h.handleAuthSuccess(c, response.Token, http.StatusOK)
}

func (h *AuthWEBHandler) Register(c echo.Context) error {
	req := dto.RegisterRequest{
		FirstName: strings.TrimSpace(c.FormValue("firstname")),
		LastName:  strings.TrimSpace(c.FormValue("lastname")),
		Email:     strings.ToLower(strings.TrimSpace(c.FormValue("email"))),
		Password:  c.FormValue("password"),
	}

	if err := h.validator.Struct(req); err != nil {
		return h.renderError(c, http.StatusBadRequest, "Invalid email or password")
	}

	response, err := h.authService.Register(c.Request().Context(), &req)
	if err != nil {
		statusCode, message := h.categorizeAuthError(err)
		return h.renderError(c, statusCode, message)
	}

	return h.handleAuthSuccess(c, response.Token, http.StatusCreated)
}

func (h *AuthWEBHandler) setAuthCokie(c echo.Context, token string) {
	cookie := &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // make it true on https
		SameSite: http.SameSiteStrictMode,
		MaxAge:   360000,
	}
	c.SetCookie(cookie)
}

func (h *AuthWEBHandler) renderError(c echo.Context, statusCode int, message string) error {
	c.Response().Header().Set("Content-Type", "text/html")
	c.Response().WriteHeader(statusCode)
	return errors.AuthError(message).Render(c.Request().Context(), c.Response().Writer)
}

func (h *AuthWEBHandler) handleAuthSuccess(c echo.Context, token string, statusCode int) error {
	h.setAuthCokie(c, token)
	c.Response().Header().Set("HX-Location", "/dashboard")
	return c.NoContent(statusCode)
}

func (h *AuthWEBHandler) categorizeAuthError(err error) (int, string) {
	errStr := err.Error()

	switch {
	case strings.Contains(errStr, "invalid credentials") || strings.Contains(errStr, "user not found"):
		return http.StatusUnauthorized, "Invalid email or password"
	case strings.Contains(errStr, "already exists"):
		return http.StatusConflict, "User with this email already exists"
	default:
		return http.StatusInternalServerError, "An error occurred. Please try again"
	}
}
