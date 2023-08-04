package main

import (
	"os"

	"github.com/alekslesik/link-shorterer/internal/config"
	"golang.org/x/exp/slog"
)

const (
    envLocal = "local"
    envDev   = "dev"
    envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)

	log = log.With(slog.String("env", cfg.Env)) // к каждому сообщению будет добавляться поле с информацией о текущем окружении

	log.Info("initializing server", slog.String("address", cfg.Address))
	log.Debug("logger debug mode enabled")

}

func setupLogger(env string) *slog.Logger {
    var log *slog.Logger

    switch env {
    case envLocal:
        log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
    case envDev:
        log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
    case envProd:
        log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
    }

    return log
}