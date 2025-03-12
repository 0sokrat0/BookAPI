package main

import (
	"context"

	"github.com/0sokrat0/BookAPI/internal/application/http"
	"github.com/0sokrat0/BookAPI/internal/config"
	genid "github.com/0sokrat0/BookAPI/pkg/GenID"
	"github.com/0sokrat0/BookAPI/pkg/db/postgres"
	"github.com/0sokrat0/BookAPI/pkg/logger"

	"os"
	"os/signal"
	"syscall"
	"time"
)

var idCounter = genid.NewCounter(0)

// @title Book API
// @version 1.0
// @host localhost:8080
// @BasePath /
func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := config.GetConfig()
	if cfg == nil {
		panic("Configuration is nil")
	}

	lg := logger.NewLogger(cfg)
	defer lg.Sync()
	ctx = logger.WithLogger(ctx, lg)

	pool, err := postgres.NewPG(ctx, cfg)
	if err != nil {
		lg.Fatalf("Error connecting to PostgreSQL: %v", err)
	}
	defer pool.Close()

	server := http.NewServer(ctx, cfg, pool, idCounter)

	go func() {
		if err := server.Start(); err != nil {
			lg.Fatalf("Error starting server: %v", err)
		}
	}()

	gracefulShutdown(ctx, server, lg)
}

func gracefulShutdown(ctx context.Context, server *http.Server, lg *logger.Logger) {
	<-ctx.Done()

	lg.Info("Сигнал завершения получен, начинается graceful shutdown...")

	// Используем ctx, который содержит логгер, вместо context.Background()
	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		lg.Errorf("Ошибка при завершении работы сервера: %v", err)
	} else {
		lg.Info("Сервер успешно завершил работу")
	}
}
