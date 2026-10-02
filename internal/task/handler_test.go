package task

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type mockTaskService struct {
	createFunc  func(CreateTaskRequest) (Task, error)
	getByIDFunc func(string) (Task, error)
}

func (m mockTaskService) Create(request CreateTaskRequest) (Task, error) {
	return m.createFunc(request)
}

func (m mockTaskService) GetByID(id string) (Task, error) {
	return m.getByIDFunc(id)
}

func TestTaskHandler_Post(t *testing.T) {
	createdAt := time.Now()
	updatedAt := createdAt

	expectedTask := Task{
		ID:        "01a0fb8d-7872-70ec-8000-baabc6874a55",
		Type:      TaskType("email"),
		Payload:   json.RawMessage(`{"to":"user@example.com"}`),
		Status:    TaskStatusPending,
		Attempts:  0,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	service := mockTaskService{
		createFunc: func(request CreateTaskRequest) (Task, error) {
			return expectedTask, nil
		},
	}

	handler := NewTaskHandler(service)

	request := httptest.NewRequest(
		http.MethodPost,
		"/tasks",
		strings.NewReader(`{
			"type": "email",
			"payload": {
				"to": "user@example.com"
			}
		}`),
	)

	recorder := httptest.NewRecorder()

	handler.Post(recorder, request)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, recorder.Code)
	}

	if contentType := recorder.Header().Get("Content-Type"); contentType != "application/json" {
		t.Fatalf("expected Content-Type application/json, got %q", contentType)
	}

	var response TaskResponse

	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.ID != expectedTask.ID {
		t.Errorf("expected ID %q, got %q", expectedTask.ID, response.ID)
	}

	if response.Type != expectedTask.Type {
		t.Errorf("expected type %q, got %q", expectedTask.Type, response.Type)
	}

	if string(response.Payload) != string(expectedTask.Payload) {
		t.Errorf("expected payload %s, got %s", expectedTask.Payload, response.Payload)
	}

	if response.Status != expectedTask.Status {
		t.Errorf("expected status %q, got %q", expectedTask.Status, response.Status)
	}

	if response.Attempts != expectedTask.Attempts {
		t.Errorf("expected attempts %d, got %d", expectedTask.Attempts, response.Attempts)
	}

	if !response.CreatedAt.Equal(expectedTask.CreatedAt) {
		t.Errorf("expected createdAt %v, got %v", expectedTask.CreatedAt, response.CreatedAt)
	}

	if !response.UpdatedAt.Equal(expectedTask.UpdatedAt) {
		t.Errorf("expected updatedAt %v, got %v", expectedTask.UpdatedAt, response.UpdatedAt)
	}
}

func TestTaskHandler_Post_InvalidJSON(t *testing.T) {
	service := mockTaskService{
		createFunc: func(request CreateTaskRequest) (Task, error) {
			t.Fatal("Create should not be called")
			return Task{}, nil
		},
	}

	handler := NewTaskHandler(service)

	request := httptest.NewRequest(
		http.MethodPost,
		"/tasks",
		strings.NewReader(`{"type":`),
	)

	recorder := httptest.NewRecorder()

	handler.Post(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, recorder.Code)
	}

	assertErrorResponse(t, recorder, "invalid request body")
}

func TestTaskHandler_Post_MissingType(t *testing.T) {
	service := mockTaskService{
		createFunc: func(request CreateTaskRequest) (Task, error) {
			t.Fatal("Create should not be called")
			return Task{}, nil
		},
	}

	handler := NewTaskHandler(service)

	request := httptest.NewRequest(
		http.MethodPost,
		"/tasks",
		strings.NewReader(`{"payload":{"to":"user@example.com"}}`),
	)

	recorder := httptest.NewRecorder()

	handler.Post(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, recorder.Code)
	}

	assertErrorResponse(t, recorder, "type is required")
}

func TestTaskHandler_Post_ServiceError(t *testing.T) {
	service := mockTaskService{
		createFunc: func(request CreateTaskRequest) (Task, error) {
			return Task{}, errors.New("repository error")
		},
	}

	handler := NewTaskHandler(service)

	request := httptest.NewRequest(
		http.MethodPost,
		"/tasks",
		strings.NewReader(`{"type":"email"}`),
	)

	recorder := httptest.NewRecorder()

	handler.Post(recorder, request)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusInternalServerError,
			recorder.Code,
		)
	}

	assertErrorResponse(t, recorder, "internal server error")
}

func TestTaskHandler_Get(t *testing.T) {
	expectedTask := Task{
		ID:        "01a0fb8d-7872-70ec-8000-baabc6874a55",
		Type:      TaskType("email"),
		Payload:   json.RawMessage(`{"to":"user@example.com"}`),
		Status:    TaskStatusPending,
		Attempts:  0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	service := mockTaskService{
		getByIDFunc: func(id string) (Task, error) {
			if id != expectedTask.ID {
				t.Errorf("expected ID %q, got %q", expectedTask.ID, id)
			}

			return expectedTask, nil
		},
	}

	handler := NewTaskHandler(service)

	request := httptest.NewRequest(
		http.MethodGet,
		"/tasks/"+expectedTask.ID,
		nil,
	)

	request.SetPathValue("id", expectedTask.ID)

	recorder := httptest.NewRecorder()

	handler.Get(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	if contentType := recorder.Header().Get("Content-Type"); contentType != "application/json" {
		t.Fatalf("expected Content-Type application/json, got %q", contentType)
	}

	var response TaskResponse

	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.ID != expectedTask.ID {
		t.Errorf("expected ID %q, got %q", expectedTask.ID, response.ID)
	}

	if response.Type != expectedTask.Type {
		t.Errorf("expected type %q, got %q", expectedTask.Type, response.Type)
	}

	if string(response.Payload) != string(expectedTask.Payload) {
		t.Errorf("expected payload %s, got %s", expectedTask.Payload, response.Payload)
	}

	if response.Status != expectedTask.Status {
		t.Errorf("expected status %q, got %q", expectedTask.Status, response.Status)
	}

	if response.Attempts != expectedTask.Attempts {
		t.Errorf("expected attempts %d, got %d", expectedTask.Attempts, response.Attempts)
	}
}

func TestTaskHandler_Get_NotFound(t *testing.T) {
	service := mockTaskService{
		getByIDFunc: func(id string) (Task, error) {
			return Task{}, ErrTaskNotFound
		},
	}

	handler := NewTaskHandler(service)

	request := httptest.NewRequest(
		http.MethodGet,
		"/tasks/01a0fb8d-7872-70ec-8000-baabc6874a55",
		nil,
	)

	request.SetPathValue("id", "01a0fb8d-7872-70ec-8000-baabc6874a55")

	recorder := httptest.NewRecorder()

	handler.Get(recorder, request)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, recorder.Code)
	}

	assertErrorResponse(t, recorder, "task not found")
}

func TestTaskHandler_Get_ServiceError(t *testing.T) {
	service := mockTaskService{
		getByIDFunc: func(id string) (Task, error) {
			return Task{}, errors.New("repository error")
		},
	}

	handler := NewTaskHandler(service)

	request := httptest.NewRequest(
		http.MethodGet,
		"/tasks/task-123",
		nil,
	)

	request.SetPathValue("id", "01a0fb8d-7872-70ec-8000-baabc6874a55")

	recorder := httptest.NewRecorder()

	handler.Get(recorder, request)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusInternalServerError,
			recorder.Code,
		)
	}

	assertErrorResponse(t, recorder, "internal server error")
}

func assertErrorResponse(t *testing.T, recorder *httptest.ResponseRecorder, expectedMessage string) {
	t.Helper()

	var response struct {
		Status string `json:"status"`
		Error  string `json:"error"`
	}

	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}

	if response.Status != "error" {
		t.Errorf("expected status %q, got %q", "error", response.Status)
	}

	if response.Error != expectedMessage {
		t.Errorf("expected error %q, got %q", expectedMessage, response.Error)
	}
}
