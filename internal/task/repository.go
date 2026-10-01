package task

type TaskRepository interface {
	Create(task Task) error
	GetById(id string) (Task, error)
}
