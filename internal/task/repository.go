package task

type TaskRepository interface {
	Create(task Task) error
	GetByID(id string) (Task, error)
}
