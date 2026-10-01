package task

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestTaskHandler_CreateTask(t *testing.T) {
	repository := NewMemoryRepository()
	service := NewTaskService(repository)
	handler := NewTaskHandler(service)

	body := `{
		"type": "email",
		"payload": {
			"to": "user@example.com"
		}
	}`

	request := httptest.NewRequest(
		http.MethodPost,
		"/tasks",
		strings.NewReader(body),
	)

	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, recorder.Code)
	}

	if contentType := recorder.Header().Get("Content-Type"); contentType != "application/json" {
		t.Fatalf("expected Content-Type application/json, got %q", contentType)
	}

	var response CreateTaskResponse

	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.ID == "" {
		t.Error("expected task ID to be present")
	}

	if response.Type != TaskType("email") {
		t.Errorf("expected type %q, got %q", "email", response.Type)
	}

	if response.Status != TaskStatusPending {
		t.Errorf("expected status %q, got %q", TaskStatusPending, response.Status)
	}

	expectedPayload := `{"to":"user@example.com"}`

	if string(response.Payload) != expectedPayload {
		t.Errorf("expected payload %s, got %s", expectedPayload, response.Payload)
	}
}

func TestTaskHandler_InvalidJSON(t *testing.T) {
	repository := NewMemoryRepository()
	service := NewTaskService(repository)
	handler := NewTaskHandler(service)

	request := httptest.NewRequest(
		http.MethodPost,
		"/tasks",
		strings.NewReader(`{"type":`),
	)

	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, recorder.Code)
	}

	var response struct {
		Status string `json:"status"`
		Error  string `json:"error"`
	}

	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Status != "error" {
		t.Errorf("expected status %q, got %q", "error", response.Status)
	}

	if response.Error != "invalid request body" {
		t.Errorf("expected error %q, got %q", "invalid request body", response.Error)
	}
}

func TestTaskHandler_MissingType(t *testing.T) {
	repository := NewMemoryRepository()
	service := NewTaskService(repository)
	handler := NewTaskHandler(service)

	request := httptest.NewRequest(
		http.MethodPost,
		"/tasks",
		strings.NewReader(`{"payload":{"to":"user@example.com"}}`),
	)

	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, recorder.Code)
	}

	var response struct {
		Status string `json:"status"`
		Error  string `json:"error"`
	}

	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Status != "error" {
		t.Errorf("expected status %q, got %q", "error", response.Status)
	}

	if response.Error != "type is required" {
		t.Errorf("expected error %q, got %q", "type is required", response.Error)
	}
}

func TestTaskHandler_CreateTask_GeneratesUniqueIDs(t *testing.T) {
	repository := NewMemoryRepository()
	service := NewTaskService(repository)
	handler := NewTaskHandler(service)

	body := `{"type":"email","payload":{"to":"user@example.com"}}`

	request1 := httptest.NewRequest(
		http.MethodPost,
		"/tasks",
		strings.NewReader(body),
	)

	recorder1 := httptest.NewRecorder()

	handler.ServeHTTP(recorder1, request1)

	if recorder1.Code != http.StatusCreated {
		t.Fatalf("expected first status %d, got %d", http.StatusCreated, recorder1.Code)
	}

	var response1 CreateTaskResponse

	if err := json.NewDecoder(recorder1.Body).Decode(&response1); err != nil {
		t.Fatalf("failed to decode first response: %v", err)
	}

	request2 := httptest.NewRequest(
		http.MethodPost,
		"/tasks",
		strings.NewReader(body),
	)

	recorder2 := httptest.NewRecorder()

	handler.ServeHTTP(recorder2, request2)

	if recorder2.Code != http.StatusCreated {
		t.Fatalf("expected second status %d, got %d", http.StatusCreated, recorder2.Code)
	}

	var response2 CreateTaskResponse

	if err := json.NewDecoder(recorder2.Body).Decode(&response2); err != nil {
		t.Fatalf("failed to decode second response: %v", err)
	}

	if response1.ID == response2.ID {
		t.Errorf("expected unique IDs, got the same ID %q", response1.ID)
	}
}
