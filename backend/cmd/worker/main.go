package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	broadcastmod "github.com/beatfraps/wa-dashboard/backend/internal/modules/broadcast"
	"github.com/beatfraps/wa-dashboard/backend/internal/shared/config"
	"github.com/beatfraps/wa-dashboard/backend/internal/shared/db"
	"github.com/beatfraps/wa-dashboard/backend/internal/shared/logx"
	"github.com/beatfraps/wa-dashboard/backend/internal/shared/queue"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "config error: %v\n", err)
		os.Exit(1)
	}

	logger := logx.New(cfg.LogLevel)

	ctx := context.Background()
	pool, err := db.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Error("failed to connect database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	qClient, err := queue.NewClient(cfg.RedisURL)
	if err != nil {
		logger.Error("failed to connect redis", "error", err)
		os.Exit(1)
	}
	defer qClient.Close()

	worker, err := queue.NewWorker(cfg.RedisURL, 10)
	if err != nil {
		logger.Error("failed to create worker", "error", err)
		os.Exit(1)
	}

	broadcastModule := broadcastmod.NewModule(pool.Pool, qClient)
	broadcastModule.RegisterWorkerHandlers(worker, logger)

	go func() {
		logger.Info("starting asynq worker")
		if err := worker.Run(); err != nil {
			logger.Error("worker error", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down worker")
	worker.Shutdown()
}
