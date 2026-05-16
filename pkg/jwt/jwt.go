package jwt

import (
	"cinema-ticketing-api/internal/config"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JwtClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateJwt(UserId uuid.UUID, Role string) (string, error) {
	expiresAt := config.GetEnv("JWT_EXPIRATION", "24h")

	duration, err := time.ParseDuration(expiresAt)
	if err != nil {
		return "", err
	}
	claims := JwtClaims{
		UserID: UserId.String(),
		Role:   Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := config.GetEnv("JWT_SECRET", "secret_key")

	return token.SignedString([]byte(secret))
}

func ValidateToken(tokenString string) (*JwtClaims, error) {
	secretKey := []byte(config.GetEnv("JWT_SECRET", "secret_key"))

	token, err := jwt.ParseWithClaims(tokenString, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JwtClaims)

	if !ok || !token.Valid {
		return nil, errors.New("Invalid token")
	}

	return claims, nil
}
