package refresh

import (
	"errors"
	"fmt"
	"os"

	"github.com/MTUCIBOY/MedodsTest/pkg/tokens"
	"github.com/golang-jwt/jwt/v5"
)

var errEmptySecretKey = errors.New("JWT_SECRET is not set")

func New(accessToken string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS512)

	// Благодаря этой связке выполняется требование о парах
	secretKey := []byte(os.Getenv("JWT_SECRET") + accessToken)
	if len(secretKey) == 0 {
		return "", errEmptySecretKey
	}

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to signed token: %w", err)
	}

	return tokenString, nil
}

func CheckToken(refreshToken, accessToken string) error {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return tokens.ErrEmptySecretKey
	}

	jwtSecret += accessToken

	parsedToken, err := jwt.Parse(refreshToken, func(_ *jwt.Token) (any, error) {
		return []byte(jwtSecret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS512.Alg()}))
	if err != nil {
		return fmt.Errorf("failed to parse token: %w", err)
	}

	if !parsedToken.Valid {
		return tokens.ErrInvalidToken
	}

	return nil
}
