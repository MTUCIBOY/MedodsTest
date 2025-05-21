package storage

import (
	"context"
	"errors"
)

type DB interface {
	Connect(ctx context.Context, dsn string) error
	Auth(ctx context.Context, email, password string) (bool, error)
	AddUser(ctx context.Context, email, password string) error
	AddRefreshToken(ctx context.Context, email, token, uuidToken string) error
	UserUUID(ctx context.Context, email string) (string, error)
	IsActiveRefresh(ctx context.Context, token, uuidToken string) (bool, error)
	DeactivateRefreshToken(ctx context.Context, uuidToken string) error
	Close()
}

var (
	ErrEmailExist     = errors.New("email exists")
	ErrEmailNotFound  = errors.New("email not found")
	ErrWrongPassword  = errors.New("wrong password")
	ErrNotMatchHashes = errors.New("the hashes don't match")
)

const (
	MaxLenRefreshHash = 72
)

const (
	InitDBQuery = `
		CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

		CREATE TABLE IF NOT EXISTS users (
			uuid UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			email TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT NOW()
		);

		CREATE TABLE IF NOT EXISTS refresh_hashes (
				uuid UUID PRIMARY KEY,
				user_uuid UUID REFERENCES users(uuid),
				token_hash TEXT NOT NULL,
				is_active BOOL DEFAULT TRUE NOT NULL,
				created_at TIMESTAMP DEFAULT NOW() NOT NULL 
		);
	`

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
		INSERT INTO refresh_hashes (uuid, user_uuid, token_hash)
		VALUES ($1, $2, $3);
	`

	UserUUIDQuery = `
		SELECT uuid
		FROM users
		WHERE email = $1;
	`

	//nolint:gosec
	RefreshTokenQuery = `
		SELECT token_hash, is_active
		FROM refresh_hashes
		WHERE uuid = $1;
	`

	//nolint:gosec
	DeactivateRefreshTokenQuery = `
		UPDATE refresh_hashes
		SET is_active = false
		WHERE uuid = $1;
	`
)
