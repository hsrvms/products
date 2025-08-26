package services

import (
	"fmt"
	"go-starter/internal/modules/auth/models"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	jwt.RegisteredClaims
}

type JWTService struct {
	secretKey []byte
	issuer    string
}

func NewJWTService() *JWTService {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "hard-coded-secret"
	}

	issuer := os.Getenv("JWT_ISSUER")
	if issuer == "" {
		issuer = "go-starter"
	}

	return &JWTService{
		secretKey: []byte(secret),
		issuer:    issuer,
	}
}

func (j *JWTService) GenerateToken(user *models.User) (string, error) {
	expirationTimeStr := os.Getenv("JWT_EXPIRATION_HOURS")
	expirationHours := 24 // default 24 hours

	if expirationTimeStr != "" {
		if hours, err := strconv.Atoi(expirationTimeStr); err == nil {
			expirationHours = hours
		}
	}

	expirationTime := time.Now().Add(time.Duration(expirationHours) * time.Hour)

	claims := &Claims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    j.issuer,
			Subject:   user.ID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(j.secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (j *JWTService) ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

func GetUserIDFromContext(c echo.Context) (uuid.UUID, error) {
	user := c.Get("user")
	if user == nil {
		return uuid.Nil, fmt.Errorf("user not found in context")
	}

	claims, ok := user.(*Claims)
	if !ok {
		return uuid.Nil, fmt.Errorf("invalid user claims in context")
	}

	return claims.UserID, nil
}
