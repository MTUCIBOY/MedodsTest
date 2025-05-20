package tokens

import "errors"

var (
	ErrEmptySecretKey = errors.New("JWT_SECRET is not set")
	ErrInvalidClaims  = errors.New("not valid claims")
)
