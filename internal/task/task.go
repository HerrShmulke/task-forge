package task

import (
	"encoding/json"
	"errors"
	"time"
)

type TaskType string
type TaskStatus string

const (
	TaskStatusPending   TaskStatus = "pending"
	TaskStatusRunning   TaskStatus = "running"
	TaskStatusCompleted TaskStatus = "completed"
	TaskStatusFailed    TaskStatus = "failed"
)

var ErrTaskNotFound = errors.New("task not found")
var ErrInvalidTaskId = errors.New("invalid task id")

type CreateTaskRequest struct {
	Type    TaskType
	Payload json.RawMessage
}

type TaskResponse struct {
	ID        string          `json:"id"`
	Type      TaskType        `json:"type"`
	Payload   json.RawMessage `json:"payload"`
	Status    TaskStatus      `json:"status"`
	Attempts  int             `json:"attempts"`
	CreatedAt time.Time       `json:"createdAt"`
	UpdatedAt time.Time       `json:"updatedAt"`
}

type Task struct {
	ID        string
	Type      TaskType
	Payload   json.RawMessage
	Status    TaskStatus
	Attempts  int
	CreatedAt time.Time
	UpdatedAt time.Time
}
