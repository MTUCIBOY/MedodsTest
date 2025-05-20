package refresh

import (
	"errors"
	"fmt"
	"os"

	"github.com/MTUCIBOY/MedodsTest/pkg/tokens"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var errEmptySecretKey = errors.New("JWT_SECRET is not set")

func New() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.RegisteredClaims{
		ID: uuid.New().String(),
	})

	secretKey := []byte(os.Getenv("JWT_SECRET"))
	if len(secretKey) == 0 {
		return "", errEmptySecretKey
	}

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to signed token: %w", err)
	}

	return tokenString, nil
}

func Check(token string) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return "", tokens.ErrEmptySecretKey
	}

	parsedToken, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(_ *jwt.Token) (any, error) {
		return []byte(jwtSecret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS512.Alg()}))
	if err != nil {
		return "", fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := parsedToken.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return "", tokens.ErrInvalidClaims
	}

	return claims.ID, nil
}
