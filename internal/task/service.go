package task

import "github.com/google/uuid"

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

	task := Task{
		ID:      id.String(),
		Type:    request.Type,
		Payload: request.Payload,
		Status:  TaskStatusPending,
	}

	if err := s.repository.Create(task); err != nil {
		return Task{}, err
	}

	return task, nil
}

func (s *TaskService) GetById(id string) (Task, error) {
	task, err := s.repository.GetById(id)

	if err != nil {
		return Task{}, err
	}

	return task, nil
}
