package storage

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/ichbingautam/distributed-task-scheduler/internal/core"
)

type RedisStore struct {
	client *redis.Client
	prefix string
}

func NewRedisStore(address, password string, db int) *RedisStore {
	return &RedisStore{
		client: redis.NewClient(&redis.Options{
			Addr:     address,
			Password: password,
			DB:       db,
		}),
		prefix: "tasks",
	}
}

var _ Store = (*RedisStore)(nil)

func (s *RedisStore) SaveTask(ctx context.Context, task *core.Task) error {
	data, err := json.Marshal(task)
	if err != nil {
		return err
	}
	return s.client.ZAdd(ctx, s.prefix+":queue", &redis.Z{
		Score:  float64(task.ScheduledAt.UnixNano()),
		Member: data,
	}).Err()
}

func (s *RedisStore) LoadTasks(ctx context.Context) ([]*core.Task, error) {
	data, err := s.client.ZRangeByScore(ctx, s.prefix+":queue", &redis.ZRangeBy{
		Min: "0",
		Max: "+inf",
	}).Result()

	if err != nil {
		return nil, err
	}

	var tasks []*core.Task
	for _, d := range data {
		var task core.Task
		if err := json.Unmarshal([]byte(d), &task); err != nil {
			continue
		}
		tasks = append(tasks, &task)
	}
	return tasks, nil
}