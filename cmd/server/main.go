package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"taskforge/internal/config"
	"taskforge/internal/health"
	"taskforge/internal/task"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Printf("failed to load .env: %v\n", err)
		return
	}

	cfg, err := config.LoadFromEnv()
	if err != nil {

		fmt.Printf("config load error: %v\n", err)
		return
	}

	mux := http.NewServeMux()

	healthHandler := health.HealthHandler{}
	mux.Handle("GET /health", healthHandler)

	taskRepository := task.NewMemoryRepository()
	taskService := task.NewTaskService(taskRepository)
	taskHandler := task.NewTaskHandler(taskService)
	mux.HandleFunc("POST /tasks", taskHandler.Post)
	mux.HandleFunc("GET /tasks/{id}", taskHandler.Get)

	server := &http.Server{
		Addr:    ":" + strconv.Itoa(cfg.HTTPPort),
		Handler: mux,
	}

	errCh := make(chan error)

	go func() {
		err := server.ListenAndServe()

		if !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	signals := make(chan os.Signal)
	defer signal.Stop(signals)

	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	select {
	case sig := <-signals:
		fmt.Printf("received signal: %v\n", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := shutdown(server, ctx); err != nil {
			fmt.Printf("server shutdown error: %v\n", err)
		}
	case err := <-errCh:
		fmt.Printf("server shutdown error: %v\n", err)
	}
}
