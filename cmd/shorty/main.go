package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/kiryu-dev/shorty/internal/config"
	"github.com/kiryu-dev/shorty/internal/http"
	"github.com/kiryu-dev/shorty/internal/service"
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
	logger := setupLogger(config.Env)
	logger.Info("shorty is here 😼", slog.String("env", config.Env))
	storage, err := postgres.New(&config.DB)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	defer storage.Close()
	var (
		urlService = service.NewShortener(storage)
		httpServer = http.New(&config.HTTPServer, urlService)
	)
	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			logger.Error(err.Error())
			return
		}
	}()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	<-sigChan
	logger.Info("bye dude 😼")
	if err := httpServer.Shutdown(context.Background()); err != nil {
		logger.Error(err.Error())
		return
	}
}

func setupLogger(env string) *slog.Logger {
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
