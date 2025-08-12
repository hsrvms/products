package middlewares

import (
	"net/http"
	"products/internal/modules/auth/services"
	"strings"

	"github.com/labstack/echo/v4"
)

func JWTMiddleware(jwtService *services.JWTService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get token from header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Authorization header required",
				})
			}

			// Check if header starts with "Bearer "
			if !strings.HasPrefix(authHeader, "Bearer ") {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Invalid authorization header format",
				})
			}

			// Extract token
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			// Validate token
			claims, err := jwtService.ValidateToken(tokenString)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Invalid or expired token",
				})
			}

			// Set user in context
			c.Set("user", claims)

			return next(c)
		}
	}
}
