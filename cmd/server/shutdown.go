package main

import (
	"context"
	"net/http"
)

func shutdown(server *http.Server, ctx context.Context) error {
	return server.Shutdown(ctx)
}
