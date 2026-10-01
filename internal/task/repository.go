package task

type TaskRepository interface {
	Create(task Task) error
}
