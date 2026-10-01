package health

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	healthHandler := HealthHandler{}

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/health", nil)

	healthHandler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Errorf("expected status %v, got %v", http.StatusOK, recorder.Code)
	}

	const expectedContentType = "application/json"

	contentType := recorder.Header().Get("Content-Type")

	if contentType != expectedContentType {
		t.Errorf("expected Content-Type %v, got %v", expectedContentType, contentType)
	}

	var healthResponse HealthResponse

	jsonData := recorder.Body.String()
	err := json.Unmarshal([]byte(jsonData), &healthResponse)

	if err != nil {
		t.Fatalf("failed to deserialize json")
	}

	const expectedStatus = "ok"

	if healthResponse.Status != expectedStatus {
		t.Errorf("expected status %v, got %v", expectedStatus, healthResponse.Status)
	}
}
