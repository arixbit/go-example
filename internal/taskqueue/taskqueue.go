// Package taskqueue defines a typed async queue boundary for application tasks.
package taskqueue

import (
	"context"
	"errors"

	"github.com/hibiken/asynq"
)

var (
	// ErrQueueUnavailable is returned when a queue wrapper has no underlying client.
	ErrQueueUnavailable = errors.New("taskqueue: queue unavailable")
	// ErrNilTask is returned when callers try to enqueue a nil task.
	ErrNilTask = errors.New("taskqueue: nil task")
)

// Queue publishes background tasks.
type Queue struct {
	client *asynq.Client
}

// NewQueue wraps an asynq client. It returns nil when client is nil.
func NewQueue(client *asynq.Client) *Queue {
	if client == nil {
		return nil
	}
	return &Queue{client: client}
}

// Available reports whether the queue has an underlying client.
func (q *Queue) Available() bool {
	return q != nil && q.client != nil
}

// Enqueue publishes a task to the queue.
func (q *Queue) Enqueue(ctx context.Context, t *asynq.Task, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	if q == nil || q.client == nil {
		return nil, ErrQueueUnavailable
	}
	if t == nil {
		return nil, ErrNilTask
	}
	if ctx == nil {
		ctx = context.Background()
	}
	return q.client.EnqueueContext(ctx, t, opts...)
}
