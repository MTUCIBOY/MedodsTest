package psql

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
)

type Storage struct {
	db     *pgx.Conn
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

	conConfig, err := pgx.ParseConfig(dsn)
	if err != nil {
		log.Error("failed to parse config", slog.String("err", err.Error()))

		return fmt.Errorf("failed to parse config: %w", err)
	}

	conn, err := pgx.ConnectConfig(ctx, conConfig)
	if err != nil {
		log.Error("failed to connect to DB", slog.String("err", err.Error()))

		return fmt.Errorf("failed to connect to DB: %w", err)
	}

	s.db = conn

	log.Info("Connected to database")

	return nil
}

func (s *Storage) Close(ctx context.Context) {
	const fn = "psql.Storage.Close"
	log := s.logger.With(slog.String("fn", fn))

	if err := s.db.Close(ctx); err != nil {
		log.Error("failed to close DB", slog.String("err", err.Error()))
	}

	log.Info("Connection to database is closed")
}
