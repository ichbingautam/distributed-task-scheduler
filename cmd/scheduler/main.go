package main

import (
	"context"
	"log"
	"github.com/ichbingautam/distributed-task-scheduler/config"
	"github.com/ichbingautam/distributed-task-scheduler/internal/executor"
	"github.com/ichbingautam/distributed-task-scheduler/internal/metrics"
	"github.com/ichbingautam/distributed-task-scheduler/internal/scheduler"
	"github.com/ichbingautam/distributed-task-scheduler/internal/storage"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Initialize metrics
	metrics.Init(cfg.MetricsPort)

	// Initialize storage
	store := storage.NewRedisStore(
		cfg.Redis.Address,
		cfg.Redis.Password,
		cfg.Redis.DB,
	)

	// Create scheduler and executor
	sched := scheduler.NewScheduler(store)      // Now matches signature
	exec := executor.NewExecutor(cfg, sched)    // Updated executor constructor

	// Start the system
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	exec.Start(ctx)

	// Block until shutdown signal
	<-ctx.Done()
	exec.Stop()
}