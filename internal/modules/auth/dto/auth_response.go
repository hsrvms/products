package dto

import "products/internal/modules/auth/models"

type AuthResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}
