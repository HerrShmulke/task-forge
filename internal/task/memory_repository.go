package task

import (
	"fmt"
	"sync"
)

type MemoryRepository struct {
	tasks map[string]Task
	mu    sync.Mutex
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		tasks: make(map[string]Task),
	}
}

func (r *MemoryRepository) Create(task Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tasks[task.ID]; exists {
		return fmt.Errorf("task with ID %s already exists", task.ID)
	}

	r.tasks[task.ID] = task

	return nil
}

func (r *MemoryRepository) GetByID(id string) (Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	task, exists := r.tasks[id]

	if !exists {
		return Task{}, ErrTaskNotFound
	}

	return task, nil
}

var _ TaskRepository = (*MemoryRepository)(nil)
