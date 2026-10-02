package task

import (
	"time"

	"github.com/google/uuid"
)

type TaskService struct {
	repository TaskRepository
}

func NewTaskService(repository TaskRepository) *TaskService {
	return &TaskService{
		repository: repository,
	}
}

func (s *TaskService) Create(request CreateTaskRequest) (Task, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return Task{}, err
	}

	now := time.Now()

	task := Task{
		ID:        id.String(),
		Type:      request.Type,
		Payload:   request.Payload,
		Status:    TaskStatusPending,
		Attempts:  0,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.repository.Create(task); err != nil {
		return Task{}, err
	}

	return task, nil
}

func (s *TaskService) GetByID(id string) (Task, error) {
	err := ParseTaskID(id)
	if err != nil {
		return Task{}, ErrInvalidTaskId
	}

	task, err := s.repository.GetByID(id)

	if err != nil {
		return Task{}, err
	}

	return task, nil
}

func ParseTaskID(value string) error {
	return uuid.Validate(value)
}
