package main

import (
	"flag"
	"log"
	"os"

	"github.com/kiryu-dev/shorty/internal/config"
	"github.com/kiryu-dev/shorty/internal/storage/postgres"
	"golang.org/x/exp/slog"
)

const (
	envLocal = "local"
	envDev   = "dev"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "./configs/local.yaml", "config file path")
	flag.Parse()
	config, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}
	logger := initLogger(config.Env)
	logger.Info("shorty is here 😼", slog.String("env", config.Env))
	storage, err := postgres.New(&config.DB)
	if err != nil {
		logger.Debug(err.Error())
	}
	defer storage.Close()
}

func initLogger(env string) *slog.Logger {
	switch env {
	case envLocal:
		return slog.New(
			slog.NewTextHandler(os.Stdin, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		return slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	}
	return nil
}
