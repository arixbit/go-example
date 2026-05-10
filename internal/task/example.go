package task

import (
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
)

const (
	// TypeExampleTask identifies the example async task.
	TypeExampleTask = "example:run"
)

// ExamplePayload is the payload for the example task.
type ExamplePayload struct {
	Name    string `json:"name"`
	TraceID string `json:"trace_id,omitempty"`
}

// NewExampleTask creates a new example task for async processing.
func NewExampleTask(name string) (*asynq.Task, error) {
	payload, err := json.Marshal(ExamplePayload{Name: name})
	if err != nil {
		return nil, fmt.Errorf("marshal example payload: %w", err)
	}
	return asynq.NewTask(TypeExampleTask, payload, asynq.MaxRetry(5)), nil
}
