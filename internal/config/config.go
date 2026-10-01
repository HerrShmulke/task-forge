package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	HTTPPort    int
	DatabaseURL string
	WorkerCount int
}

func LoadFromEnv() (Config, error) {
	httpPortStr := os.Getenv("HTTP_PORT")
	if httpPortStr == "" {
		return Config{}, errors.New("HTTP_PORT is empty")
	}

	httpPort, err := strconv.Atoi(httpPortStr)
	if err != nil {
		return Config{}, fmt.Errorf("failed to parse HTTP_PORT: %w", err)
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return Config{}, errors.New("DATABASE_URL is empty")
	}

	workerCountStr := os.Getenv("WORKER_COUNT")
	if workerCountStr == "" {
		return Config{}, errors.New("WORKER_COUNT is empty")
	}

	workerCount, err := strconv.Atoi(workerCountStr)
	if err != nil {
		return Config{}, fmt.Errorf("failed to parse WORKER_COUNT: %w", err)
	}

	if httpPort <= 0 || httpPort > 65535 {
		return Config{}, errors.New("invalid HTTP_PORT")
	}

	if workerCount <= 0 {
		return Config{}, errors.New("invalid WORKER_COUNT")
	}

	return Config{
		HTTPPort:    httpPort,
		DatabaseURL: databaseURL,
		WorkerCount: workerCount,
	}, nil
}
