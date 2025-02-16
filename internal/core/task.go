package core

import (
	"time"
	"github.com/cenkalti/backoff/v4"
)

type Task struct {
	ID          string
	Execute     func() error
	ScheduledAt time.Time
	RetryPolicy RetryPolicy
	Attempts    int
}

type RetryPolicy struct {
	MaxAttempts int
	Backoff     backoff.BackOff
}

func NewExponentialRetryPolicy(maxAttempts int) RetryPolicy {
	b := backoff.NewExponentialBackOff()
	b.InitialInterval = 1 * time.Second
	b.MaxInterval = 1 * time.Minute
	b.Multiplier = 2
	b.Reset()

	return RetryPolicy{
		MaxAttempts: maxAttempts,
		Backoff:     backoff.WithMaxRetries(b, uint64(maxAttempts)),
	}
}

type RetryableError struct {
	Err error
}

func (e RetryableError) Error() string { return e.Err.Error() }