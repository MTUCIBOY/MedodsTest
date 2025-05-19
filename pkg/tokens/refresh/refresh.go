package refresh

import (
	"errors"
	"fmt"
	"os"

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
