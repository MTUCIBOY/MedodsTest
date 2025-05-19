package access

import (
	"fmt"
	"os"
	"time"

	"github.com/MTUCIBOY/MedodsTest/pkg/tokens"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserAgent string
	jwt.RegisteredClaims
}

func New(email, userAgent string, ttl time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, Claims{
		UserAgent: userAgent,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   email,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	})

	secretKey := []byte(os.Getenv("JWT_SECRET"))
	if len(secretKey) == 0 {
		return "", tokens.ErrEmptySecretKey
	}

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to signed token: %w", err)
	}

	return tokenString, nil
}

func CheckToken(token string) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return "", tokens.ErrEmptySecretKey
	}

	parsedToken, err := jwt.ParseWithClaims(token, &Claims{}, func(_ *jwt.Token) (any, error) {
		return []byte(jwtSecret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS512.Alg()}))
	if err != nil {
		return "", fmt.Errorf("failed to parse token: %w", err)
	}

	if !parsedToken.Valid {
		return "", tokens.ErrInvalidToken
	}

	claims, ok := parsedToken.Claims.(*Claims)
	if !ok {
		return "", tokens.ErrInvalidClaims
	}

	return claims.UserAgent, nil
}
