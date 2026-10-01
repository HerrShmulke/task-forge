package task

import (
	"encoding/json"
	"net/http"
	"taskforge/internal/api"
)

type TaskHandler struct {
	service *TaskService
}

func NewTaskHandler(service *TaskService) *TaskHandler {
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

	json.NewEncoder(w).Encode(toCreateTaskResponse(newTask))
}

func (h TaskHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	task, err := h.service.GetById(id)
	if err != nil {
		// ...
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(toCreateTaskResponse(task))
}

func toCreateTaskResponse(task Task) CreateTaskResponse {
	return CreateTaskResponse{
		ID:      task.ID,
		Type:    task.Type,
		Payload: task.Payload,
		Status:  task.Status,
	}
}
