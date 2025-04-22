package handlers

import (
	"log/slog"
)

type HandlerOptions struct {
	Log *slog.Logger
}

type Resource struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
