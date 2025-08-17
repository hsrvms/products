package dto

import "go-starter/internal/modules/auth/models"

type AuthResponse struct {
	Success bool        `json:"success"`
	Token   string      `json:"token"`
	User    models.User `json:"user"`
}
