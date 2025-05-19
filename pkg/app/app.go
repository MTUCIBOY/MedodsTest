package app

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/MTUCIBOY/MedodsTest/pkg/config"
	"github.com/MTUCIBOY/MedodsTest/pkg/router"
	"github.com/MTUCIBOY/MedodsTest/pkg/storage"
)

type App struct {
	logger *slog.Logger
	server *http.Server
	db     storage.DB
}

func New(cfg *config.Config, log *slog.Logger, db storage.DB) *App {
	srv := &http.Server{
		Addr:         cfg.Address,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,

		Handler: router.New(cfg, log, db),
	}

	return &App{
		logger: log.With(slog.String("Server", cfg.Address)),
		server: srv,
		db:     db,
	}
}

func (a *App) Run() {
	const fn = "app.Run"
	log := a.logger.With(
		slog.String("fn", fn),
	)

	log.Info("Starting application...")

	err := a.db.Connect(context.TODO(), os.Getenv("STORAGE_DSN"))
	if err != nil {
		log.Error(
			"failed to connet to db",
			slog.String("err", err.Error()),
		)

		panic(err)
	}

	go gracefulShutdown(a)

	if err := a.server.ListenAndServe(); err != nil {
		log.Warn(
			"Error from ListenAndServe",
			slog.String("err", err.Error()),
		)
	}

	log.Info("Application is stopped")
}

func (a *App) Stop() {
	const fn = "app.Stop"
	log := a.logger.With(
		slog.String("fn", fn),
	)

	log.Info("Stopping application...")

	if err := a.server.Shutdown(context.TODO()); err != nil {
		log.Warn(
			"Error from server shutdown",
			slog.String("err", err.Error()),
		)
	}

	a.db.Close()
}

func gracefulShutdown(app *App) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c

	app.Stop()
}
