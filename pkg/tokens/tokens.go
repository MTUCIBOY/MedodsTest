package tokens

import "errors"

var (
	ErrEmptySecretKey = errors.New("JWT_SECRET is not set")
	ErrInvalidToken   = errors.New("not valid token")
	ErrInvalidClaims  = errors.New("not valid claims")
)
