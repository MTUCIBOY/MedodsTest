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
	UserIP    string
	jwt.RegisteredClaims
}

func New(email, userAgent, userIP, refreshToken string, ttl time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, Claims{
		UserAgent: userAgent,
		UserIP:    userIP,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   email,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	})

	// Благодаря этой связке удовлетворяется требование о парах
	secretKey := []byte(os.Getenv("JWT_SECRET") + refreshToken)
	if len(secretKey) == 0 {
		return "", tokens.ErrEmptySecretKey
	}

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to signed token: %w", err)
	}

	return tokenString, nil
}

func Check(accessToken, refreshToken string) (*Claims, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, tokens.ErrEmptySecretKey
	}

	parsedToken, err := jwt.ParseWithClaims(accessToken, &Claims{}, func(_ *jwt.Token) (any, error) {
		return []byte(jwtSecret + refreshToken), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS512.Alg()}))
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	jwtclaims, ok := parsedToken.Claims.(*Claims)
	if !ok {
		return nil, tokens.ErrInvalidClaims
	}

	return jwtclaims, nil
}

func CheckWithoutClaims(accessToken, refreshToken string) (*Claims, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, tokens.ErrEmptySecretKey
	}

	parsedToken, err := jwt.ParseWithClaims(
		accessToken,
		&Claims{},
		func(_ *jwt.Token) (any, error) {
			return []byte(jwtSecret + refreshToken), nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS512.Alg()}),
		jwt.WithoutClaimsValidation(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	jwtclaims, ok := parsedToken.Claims.(*Claims)
	if !ok {
		return nil, tokens.ErrInvalidClaims
	}

	return jwtclaims, nil
}
