package task

import (
	"encoding/json"
	"errors"
	"net/http"
	"taskforge/internal/api"
)

type taskService interface {
	Create(CreateTaskRequest) (Task, error)
	GetByID(string) (Task, error)
}

type TaskHandler struct {
	service taskService
}

func NewTaskHandler(service taskService) *TaskHandler {
	return &TaskHandler{
		service: service,
	}
}

func (h TaskHandler) Post(w http.ResponseWriter, r *http.Request) {
	req := CreateTaskRequest{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		api.WriteErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Type == "" {
		api.WriteErrorResponse(w, http.StatusBadRequest, "type is required")
		return
	}

	newTask, err := h.service.Create(req)
	if err != nil {
		api.WriteErrorResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(toTaskResponse(newTask))
}

func (h TaskHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	task, err := h.service.GetByID(id)
	if err != nil {
		if errors.Is(err, ErrInvalidTaskId) {
			api.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if errors.Is(err, ErrTaskNotFound) {
			api.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		api.WriteErrorResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(toTaskResponse(task))
}

func toTaskResponse(task Task) TaskResponse {
	return TaskResponse{
		ID:        task.ID,
		Type:      task.Type,
		Payload:   task.Payload,
		Status:    task.Status,
		Attempts:  task.Attempts,
		CreatedAt: task.CreatedAt,
		UpdatedAt: task.UpdatedAt,
	}
}
