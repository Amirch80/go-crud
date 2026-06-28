package logger

import (
	"context"
	"log/slog"
)

type MultiHandler struct {
	handlers []slog.Handler
}

func NewMultiHandler(handlers ...slog.Handler) *MultiHandler {
	return &MultiHandler{handlers: handlers}
}

func (multiHandler *MultiHandler) Enabled(context context.Context, level slog.Level) bool {
	for _, handler := range multiHandler.handlers {
		if handler.Enabled(context, level) {
			return true
		}
	}
	return false
}
