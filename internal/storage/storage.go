package storage

import (
	"context"
	"github.com/ichbingautam/distributed-task-scheduler/internal/core"
)

// Store is the persistence interface for tasks
type Store interface {
	SaveTask(ctx context.Context, task *core.Task) error
	LoadTasks(ctx context.Context) ([]*core.Task, error)
}