package storage

import (
	"context"
	"errors"
)

type DB interface {
	Connect(ctx context.Context, dsn string) error
	Auth(ctx context.Context, email, password string) (bool, error)
	AddUser(ctx context.Context, email, password string) error
	AddRefreshToken(ctx context.Context, email, token string) error
	UserUUID(ctx context.Context, email string) (string, error)
	Close()
}

var (
	ErrEmailExist    = errors.New("email exists")
	ErrEmailNotFound = errors.New("email not found")
	ErrWrongPassword = errors.New("wrong password")
)

const (
	AuthQuery = `
		SELECT password_hash 
		FROM users 
		WHERE email = $1;
	`

	NewUserQuery = `
		INSERT INTO users (email, password_hash)
		VALUES ($1, $2);
	`

	//nolint:gosec
	AddRefreshTokenQuery = `
		INSERT INTO refresh_hashes (user_uuid, token_hash)
		VALUES ($1, $2);
	`

	UserUUIDQuery = `
		SELECT uuid
		FROM users
		WHERE email = $1;
	`
)
