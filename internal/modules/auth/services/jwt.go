package services

import (
	"fmt"
	"os"
	"products/internal/modules/auth/models"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type Claims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
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
		issuer = "products"
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
			Subject:   strconv.Itoa(user.ID),
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

func GetUserIDFromContext(c echo.Context) (int, error) {
	user := c.Get("user")
	if user == nil {
		return 0, fmt.Errorf("user not found in context")
	}

	claims, ok := user.(*Claims)
	if !ok {
		return 0, fmt.Errorf("invalid user claims in context")
	}

	return claims.UserID, nil
}
