package task

import (
	"errors"
	"reflect"
	"testing"
)

func TestMemoryRepository_Create(t *testing.T) {
	repository := NewMemoryRepository()

	task := Task{
		ID:     "01a0fb8d-7872-70ec-8000-baabc6874a55",
		Type:   TaskType("email"),
		Status: TaskStatusPending,
	}

	err := repository.Create(task)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	storedTask, err := repository.GetByID(task.ID)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !reflect.DeepEqual(storedTask, task) {
		t.Errorf("expected task %+v, got %+v", task, storedTask)
	}
}

func TestMemoryRepository_Create_DuplicateID(t *testing.T) {
	repository := NewMemoryRepository()

	task := Task{
		ID:     "01a0fb8d-7872-70ec-8000-baabc6874a55",
		Type:   TaskType("email"),
		Status: TaskStatusPending,
	}

	if err := repository.Create(task); err != nil {
		t.Fatalf("failed to create first task: %v", err)
	}

	err := repository.Create(task)

	if err == nil {
		t.Fatal("expected error when creating duplicate task")
	}
}

func TestMemoryRepository_GetByID(t *testing.T) {
	repository := NewMemoryRepository()

	expectedTask := Task{
		ID:     "01a0fb8d-7872-70ec-8000-baabc6874a55",
		Type:   TaskType("email"),
		Status: TaskStatusPending,
	}

	if err := repository.Create(expectedTask); err != nil {
		t.Fatalf("failed to create task: %v", err)
	}

	task, err := repository.GetByID(expectedTask.ID)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !reflect.DeepEqual(task, expectedTask) {
		t.Errorf("expected task %+v, got %+v", expectedTask, task)
	}
}

func TestMemoryRepository_GetByID_NotFound(t *testing.T) {
	repository := NewMemoryRepository()

	_, err := repository.GetByID("01a0fb8d-7872-70ec-8000-baabc6874a55")

	if !errors.Is(err, ErrTaskNotFound) {
		t.Errorf("expected ErrTaskNotFound, got %v", err)
	}
}
