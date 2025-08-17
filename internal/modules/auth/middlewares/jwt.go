package middlewares

import (
	"go-starter/internal/modules/auth/services"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func JWTMiddleware(jwtService *services.JWTService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var tokenString string

			authHeader := c.Request().Header.Get("Authorization")
			if after, ok := strings.CutPrefix(authHeader, "Bearer "); ok {
				tokenString = after
			}

			if tokenString == "" {
				cookie, err := c.Cookie("auth_token")

				if err == nil {
					tokenString = cookie.Value
				}
			}

			if tokenString == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "no token found (header or cookie)",
				})
			}

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
