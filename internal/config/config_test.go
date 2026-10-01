package config

import (
	"errors"
	"strconv"
	"testing"
)

func TestLoadFromEnv(t *testing.T) {
	t.Setenv("HTTP_PORT", "8080")
	t.Setenv("DATABASE_URL", "http://test.ru")
	t.Setenv("WORKER_COUNT", "2")

	config, err := LoadFromEnv()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if config.HTTPPort != 8080 {
		t.Errorf("expected HTTPPort 8080, got %d", config.HTTPPort)
	}

	if config.DatabaseURL != "http://test.ru" {
		t.Errorf("expected DatabaseURL http://test.ru, got %s", config.DatabaseURL)
	}

	if config.WorkerCount != 2 {
		t.Errorf("expected WorkerCount 2, got %d", config.WorkerCount)
	}
}

func TestLoadFromEnv_HTTPPort(t *testing.T) {
	tests := []struct {
		name     string
		httpPort string
		wantErr  bool
	}{
		{
			name:     "empty",
			httpPort: "",
			wantErr:  true,
		},
		{
			name:     "invalid string",
			httpPort: "bca",
			wantErr:  true,
		},
		{
			name:     "zero",
			httpPort: "0",
			wantErr:  true,
		},
		{
			name:     "negative",
			httpPort: "-1",
			wantErr:  true,
		},
		{
			name:     "minimum valid",
			httpPort: "1",
			wantErr:  false,
		},
		{
			name:     "maximum valid",
			httpPort: "65535",
			wantErr:  false,
		},
		{
			name:     "greater than maximum",
			httpPort: "65536",
			wantErr:  true,
		},
		{
			name:     "valid",
			httpPort: "8080",
			wantErr:  false,
		},
	}

	t.Setenv("DATABASE_URL", "http://test.ru")
	t.Setenv("WORKER_COUNT", "1")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("HTTP_PORT", tt.httpPort)

			_, err := LoadFromEnv()

			if tt.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestLoadFromEnv_HTTPPortInvalidSyntax(t *testing.T) {
	t.Setenv("HTTP_PORT", "bca")
	t.Setenv("DATABASE_URL", "http://test.ru")
	t.Setenv("WORKER_COUNT", "1")

	_, err := LoadFromEnv()

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, strconv.ErrSyntax) {
		t.Errorf("expected strconv.ErrSyntax, got %v", err)
	}
}

func TestLoadFromEnv_DatabaseURLMissing(t *testing.T) {
	t.Setenv("HTTP_PORT", "8080")
	t.Setenv("DATABASE_URL", "")
	t.Setenv("WORKER_COUNT", "1")

	_, err := LoadFromEnv()

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestLoadFromEnv_WorkerCount(t *testing.T) {
	tests := []struct {
		name        string
		workerCount string
		wantErr     bool
	}{
		{
			name:        "empty",
			workerCount: "",
			wantErr:     true,
		},
		{
			name:        "invalid string",
			workerCount: "abc",
			wantErr:     true,
		},
		{
			name:        "zero",
			workerCount: "0",
			wantErr:     true,
		},
		{
			name:        "negative",
			workerCount: "-1",
			wantErr:     true,
		},
		{
			name:        "valid",
			workerCount: "1",
			wantErr:     false,
		},
	}

	t.Setenv("HTTP_PORT", "8080")
	t.Setenv("DATABASE_URL", "http://test.ru")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("WORKER_COUNT", tt.workerCount)

			_, err := LoadFromEnv()

			if tt.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
