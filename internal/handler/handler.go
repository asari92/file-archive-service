package handler

import (
	"log/slog"

	"file-archive-service/internal/service"
	"file-archive-service/pkg/config"
)

type Handler struct {
	Config  *config.Config
	Service *service.Service
	Logger  *slog.Logger
}

func NewHandler(usecases *service.Service) *Handler {
	return &Handler{Service: usecases}
}
