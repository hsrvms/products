package handlers

import (
	"net/http"
	"products/internal/modules/auth/dto"
	"products/internal/modules/auth/services"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authService *services.AuthService
	validator   *validator.Validate
}

func NewHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validator:   validator.New(),
	}
}

func (h *AuthHandler) Register(c echo.Context) error {

	var req dto.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Validation failed: " + err.Error(),
		})
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	response, err := h.authService.Register(c.Request().Context(), &req)
	if err != nil {
		// Check for specific error types
		if strings.Contains(err.Error(), "already exists") {
			return c.JSON(http.StatusConflict, map[string]string{
				"error": "User with this email already exists",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to register user",
		})
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

	return c.JSON(http.StatusCreated, response)
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Validation failed: " + err.Error(),
		})
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	response, err := h.authService.Login(c.Request().Context(), &req)
	if err != nil {
		if strings.Contains(err.Error(), "invalid credentials") || strings.Contains(err.Error(), "user not found") {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Invalid email or password",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to login",
		})
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

	return c.JSON(http.StatusOK, response)
}

func (h *AuthHandler) GetProfile(c echo.Context) error {
	userID, err := services.GetUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Unauthorized",
		})
	}

	user, err := h.authService.GetUserByID(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "User not found",
		})
	}

	return c.JSON(http.StatusOK, user)
}
