package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/beatfraps/wa-dashboard/backend/internal/app"
	"github.com/beatfraps/wa-dashboard/backend/internal/shared/config"
	"github.com/beatfraps/wa-dashboard/backend/internal/shared/logx"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "config error: %v\n", err)
		os.Exit(1)
	}

	logger := logx.New(cfg.LogLevel)
	application, err := app.New(cfg, logger)
	if err != nil {
		logger.Error("failed to start application", "error", err)
		os.Exit(1)
	}
	defer application.Close()

	addr := fmt.Sprintf(":%d", cfg.Port)
	go func() {
		logger.Info("starting api server", "addr", addr)
		if err := application.Server.Echo.Start(addr); err != nil && err != http.ErrServerClosed {
			logger.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down api server")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := application.Server.Echo.Shutdown(ctx); err != nil {
		logger.Error("shutdown error", "error", err)
	}
}
