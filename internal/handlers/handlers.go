package handlers

import (
	"log/slog"
)

type HandlerOptions struct {
	ResourceStore *[]Resource
	Log           *slog.Logger
}

type Resource struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
