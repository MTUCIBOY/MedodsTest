package psql

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/MTUCIBOY/MedodsTest/pkg/storage"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type Storage struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

func New(log *slog.Logger) *Storage {
	return &Storage{
		logger: log,
	}
}

func (s *Storage) Connect(ctx context.Context, dsn string) error {
	const fn = "psql.Storage.Connect"
	log := s.logger.With(slog.String("fn", fn))

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Error("failed to connect to DB", slog.String("err", err.Error()))

		return fmt.Errorf("failed to connect to DB: %w", err)
	}

	s.db = pool

	log.Info("Connected to database")

	return nil
}

func (s *Storage) Close() {
	s.db.Close()
	s.logger.Info("Connection to database is closed")
}

func (s *Storage) Auth(ctx context.Context, email, password string) (bool, error) {
	const fn = "psql.Storage.Auth"
	log := s.logger.With(slog.String("fn", fn))

	var passHash []byte

	err := s.db.QueryRow(ctx, storage.AuthQuery, email).Scan(&passHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Error(storage.ErrEmailNotFound.Error())

			return false, fmt.Errorf("query failed: %w", storage.ErrEmailNotFound)
		}

		log.Error("query failed", slog.String("err", err.Error()))

		return false, fmt.Errorf("query failed: %w", err)
	}

	err = bcrypt.CompareHashAndPassword(passHash, []byte(password))
	if err != nil {
		log.Error("failed compare hash", slog.String("err", err.Error()))

		return false, fmt.Errorf("failed compare hash: %w", storage.ErrWrongPassword)
	}

	return true, nil
}

func (s *Storage) AddUser(ctx context.Context, email, password string) error {
	const fn = "psql.Storage.AddUser"
	log := s.logger.With(slog.String("fn", fn))

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", slog.String("err", err.Error()))

		return fmt.Errorf("failed to generate password hash: %w", err)
	}

	_, err = s.db.Exec(ctx, storage.NewUserQuery, email, passHash)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			log.Error("not unique email", slog.String("err", err.Error()))

			return fmt.Errorf("not unique email: %w", storage.ErrEmailExist)
		}

		log.Error("query failed", slog.String("err", err.Error()))

		return fmt.Errorf("query failed: %w", err)
	}

	return nil
}

func (s *Storage) AddRefreshToken(ctx context.Context, email, token, uuidToken string) error {
	const fn = "psql.Storage.AddRefreshToken"
	log := s.logger.With(slog.String("fn", fn))

	userUUID, err := s.UserUUID(ctx, email)
	if err != nil {
		log.Error("query failed", slog.String("err", err.Error()))

		return fmt.Errorf("query failed: %w", err)
	}

	// Ограничение  длины токена. У bcrypt есть ограниечение на строку до 72 байт
	shortToken := []byte(token)[len(token)-storage.MaxLenRefreshHash:]

	tokenHash, err := bcrypt.GenerateFromPassword(shortToken, bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate token hash", slog.String("err", err.Error()))

		return fmt.Errorf("failed to generate token hash: %w", err)
	}

	_, err = s.db.Exec(ctx, storage.AddRefreshTokenQuery, uuidToken, userUUID, tokenHash)
	if err != nil {
		log.Error("query failed", slog.String("err", err.Error()))

		return fmt.Errorf("query failed: %w", err)
	}

	return nil
}

func (s *Storage) UserUUID(ctx context.Context, email string) (string, error) {
	const fn = "psql.Storage.UserUUID"
	log := s.logger.With(slog.String("fn", fn))

	var uuid string

	err := s.db.QueryRow(ctx, storage.UserUUIDQuery, email).Scan(&uuid)
	if err != nil {
		log.Error("query failed", slog.String("err", err.Error()))

		return "", fmt.Errorf("query failed: %w", err)
	}

	return uuid, nil
}

func (s *Storage) IsActiveRefresh(ctx context.Context, token, uuidToken string) (bool, error) {
	const fn = "psql.Storage.IsActiveRefresh"
	log := s.logger.With(
		slog.String("fn", fn),
	)

	var (
		refreshHash []byte
		isActive    bool
	)

	err := s.db.QueryRow(ctx, storage.RefreshTokenQuery, uuidToken).Scan(&refreshHash, &isActive)
	if err != nil {
		log.Error("query failed", slog.String("err", err.Error()))

		return isActive, fmt.Errorf("query failed: %w", err)
	}

	// Бесполезно?
	shortToken := []byte(token)[len(token)-storage.MaxLenRefreshHash:]

	err = bcrypt.CompareHashAndPassword(refreshHash, shortToken)
	if err != nil {
		log.Error(storage.ErrNotMatchHashes.Error())

		return false, storage.ErrNotMatchHashes
	}

	return isActive, nil
}

func (s *Storage) DeactivateRefreshToken(ctx context.Context, uuidToken string) error {
	const fn = "psql.Storage.DeactivateRefreshToken"
	log := s.logger.With(
		slog.String("fn", fn),
	)

	_, err := s.db.Exec(ctx, storage.DeactivateRefreshTokenQuery, uuidToken)
	if err != nil {
		log.Error("query failed", slog.String("err", err.Error()))

		return fmt.Errorf("query failed: %w", err)
	}

	return nil
}
