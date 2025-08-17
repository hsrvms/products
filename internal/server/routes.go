package server

import (
	"net/http"
	"products/internal/modules/auth"
	"products/internal/modules/dashboard"

	"github.com/labstack/echo/v4"
)

func (s *Server) initRoutes() {
	api := s.Echo.Group("/api")

	api.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	auth.RegisterRoutes(s.Echo, api, s.DB)
	dashboard.RegisterRoutes(s.Echo, api, s.DB)
}
