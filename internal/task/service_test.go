package task

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"
	"time"
)

type mockTaskRepository struct {
	createFunc  func(Task) error
	getByIDFunc func(string) (Task, error)
}

func (m mockTaskRepository) Create(task Task) error {
	return m.createFunc(task)
}

func (m mockTaskRepository) GetByID(id string) (Task, error) {
	return m.getByIDFunc(id)
}

func TestTaskService_Create(t *testing.T) {
	repository := mockTaskRepository{
		createFunc: func(task Task) error {
			return nil
		},
	}

	service := NewTaskService(repository)

	before := time.Now()

	task, err := service.Create(CreateTaskRequest{
		Type:    TaskType("email"),
		Payload: json.RawMessage(`{"to":"user@example.com"}`),
	})

	after := time.Now()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if task.ID == "" {
		t.Error("expected ID to be generated")
	}

	if task.Type != TaskType("email") {
		t.Errorf("expected type %q, got %q", "email", task.Type)
	}

	if string(task.Payload) != `{"to":"user@example.com"}` {
		t.Errorf("unexpected payload: %s", task.Payload)
	}

	if task.Status != TaskStatusPending {
		t.Errorf("expected status %q, got %q", TaskStatusPending, task.Status)
	}

	if task.Attempts != 0 {
		t.Errorf("expected attempts %d, got %d", 0, task.Attempts)
	}

	if task.CreatedAt.Before(before) || task.CreatedAt.After(after) {
		t.Errorf("CreatedAt %v is outside expected range", task.CreatedAt)
	}

	if task.UpdatedAt.Before(before) || task.UpdatedAt.After(after) {
		t.Errorf("UpdatedAt %v is outside expected range", task.UpdatedAt)
	}

	if !task.CreatedAt.Equal(task.UpdatedAt) {
		t.Errorf(
			"expected CreatedAt and UpdatedAt to be equal, got %v and %v",
			task.CreatedAt,
			task.UpdatedAt,
		)
	}
}

func TestTaskService_Create_RepositoryError(t *testing.T) {
	expectedErr := errors.New("repository error")

	repository := mockTaskRepository{
		createFunc: func(task Task) error {
			return expectedErr
		},
	}

	service := NewTaskService(repository)

	_, err := service.Create(CreateTaskRequest{
		Type:    TaskType("email"),
		Payload: json.RawMessage(`{}`),
	})

	if !errors.Is(err, expectedErr) {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}
}

func TestTaskService_GetByID(t *testing.T) {
	expectedTask := Task{
		ID:     "01a0fb8d-7872-70ec-8000-baabc6874a55",
		Type:   TaskType("email"),
		Status: TaskStatusPending,
	}

	repository := mockTaskRepository{
		getByIDFunc: func(id string) (Task, error) {
			if id != expectedTask.ID {
				t.Errorf("expected ID %q, got %q", expectedTask.ID, id)
			}

			return expectedTask, nil
		},
	}

	service := NewTaskService(repository)

	task, err := service.GetByID(expectedTask.ID)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !reflect.DeepEqual(task, expectedTask) {
		t.Errorf("expected task %+v, got %+v", expectedTask, task)
	}
}

func TestTaskService_GetByID_RepositoryError(t *testing.T) {
	expectedErr := errors.New("repository error")

	repository := mockTaskRepository{
		getByIDFunc: func(id string) (Task, error) {
			return Task{}, expectedErr
		},
	}

	service := NewTaskService(repository)

	_, err := service.GetByID("01a0fb8d-7872-70ec-8000-baabc6874a55")

	if !errors.Is(err, expectedErr) {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}
}
