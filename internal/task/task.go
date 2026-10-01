package task

import (
	"encoding/json"
	"errors"
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

type CreateTaskRequest struct {
	Type    TaskType
	Payload json.RawMessage
}

type CreateTaskResponse struct {
	ID      string          `json:"id"`
	Type    TaskType        `json:"type"`
	Payload json.RawMessage `json:"payload"`
	Status  TaskStatus      `json:"status"`
}

type Task struct {
	ID      string
	Type    TaskType
	Payload json.RawMessage
	Status  TaskStatus
}
