package executor

import (
	"context"
	"log"
	"sync"
	"time"
	"github.com/ichbingautam/distributed-task-scheduler/config"
	"github.com/ichbingautam/distributed-task-scheduler/internal/core"
	"golang.org/x/time/rate"
	"github.com/cenkalti/backoff/v4"
)

type Executor struct {
	workers    int
	scheduler SchedulerInterface
	wg         sync.WaitGroup
	limiter    *rate.Limiter
}

type SchedulerInterface interface {
	Schedule(task *core.Task)
	Start() <-chan *core.Task
	Stop()
}

func NewExecutor(cfg *config.Config, scheduler SchedulerInterface) *Executor {
	return &Executor{
		workers:    cfg.Workers,
		scheduler:  scheduler,
		limiter:    rate.NewLimiter(rate.Limit(cfg.RateLimit), cfg.RateLimit),
	}
}

func (e *Executor) Start(ctx context.Context) {
	taskChan := e.scheduler.Start()

	for i := 0; i < e.workers; i++ {
		e.wg.Add(1)
		go e.worker(ctx, taskChan)
	}
}

func (e *Executor) worker(ctx context.Context, taskChan <-chan *core.Task) {
	defer e.wg.Done()

	for {
		select {
		case task := <-taskChan:
			if task == nil {
				return
			}
			e.executeWithRetry(ctx, task)
		case <-ctx.Done():
			return
		}
	}
}

func (e *Executor) executeWithRetry(ctx context.Context, task *core.Task) {
	operation := func() error {
		err := task.Execute()
		if err != nil {
			log.Printf("Task %s attempt %d failed: %v", task.ID, task.Attempts+1, err)
			return err
		}
		return nil
	}

	notify := func(err error, delay time.Duration) {
		log.Printf("Retrying task %s in %v", task.ID, delay)
		task.ScheduledAt = time.Now().Add(delay)
		task.Attempts++
		e.scheduler.Schedule(task)
	}

	err := backoff.RetryNotify(operation, task.RetryPolicy.Backoff, notify)
	if err != nil {
		log.Printf("Task %s failed after %d attempts: %v", task.ID, task.Attempts, err)
	}
}

func (e *Executor) Stop() {
	e.scheduler.Stop()
	e.wg.Wait()
}