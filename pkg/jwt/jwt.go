package jwt

import (
	"cinema-ticketing-api/config"
	"errors"
	"time"

	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JwtClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	gojwt.RegisteredClaims
}

func GenerateJwt(userID uuid.UUID, role string, cfg config.JWTConfig) (string, error) {
	duration, err := time.ParseDuration(cfg.Expiration)
	if err != nil {
		return "", err
	}

	claims := JwtClaims{
		UserID: userID.String(),
		Role:   role,
		RegisteredClaims: gojwt.RegisteredClaims{
			ExpiresAt: gojwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  gojwt.NewNumericDate(time.Now()),
		},
	}

	token := gojwt.NewWithClaims(gojwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.Secret))
}

func ValidateToken(tokenString string, cfg config.JWTConfig) (*JwtClaims, error) {
	secretKey := []byte(cfg.Secret)

	token, err := gojwt.ParseWithClaims(tokenString, &JwtClaims{}, func(token *gojwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JwtClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
