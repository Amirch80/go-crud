package logger

import (
	"context"
	"log/slog"
)

type LevelRouter struct {
	files         map[slog.Level]slog.Handler
	fallback      slog.Handler
	console       slog.Handler
	consoleLevels map[slog.Level]bool
}

func (levelRouter *LevelRouter) Enabled(_ context.Context, _ slog.Level) bool {
	return true
}

func (levelRouter *LevelRouter) Handle(context context.Context, record slog.Record) error {
	target := levelRouter.fallback
	if handler, ok := levelRouter.files[record.Level]; ok {
		target = handler
	}
	if err := target.Handle(context, record.Clone()); err != nil {
		return err
	}
	if levelRouter.console != nil && levelRouter.consoleLevels[record.Level] {
		levelRouter.console.Handle(context, record.Clone())
	}
	return nil
}

func (levelRouter *LevelRouter) WithAttrs(attrs []slog.Attr) slog.Handler {
	files := make(map[slog.Level]slog.Handler, len(levelRouter.files))
	for level, handler := range levelRouter.files {
		files[level] = handler.WithAttrs(attrs)
	}
	newLevelRouter := &LevelRouter{
		files:         files,
		fallback:      levelRouter.fallback.WithAttrs(attrs),
		consoleLevels: levelRouter.consoleLevels,
	}
	if levelRouter.console != nil {
		newLevelRouter.console = levelRouter.console.WithAttrs(attrs)
	}
	return newLevelRouter
}

func (levelRouter *LevelRouter) WithGroup(name string) slog.Handler {
	files := make(map[slog.Level]slog.Handler, len(levelRouter.files))
	for level, handler := range levelRouter.files {
		files[level] = handler.WithGroup(name)
	}
	newLevelRouter := &LevelRouter{
		files:         files,
		fallback:      levelRouter.fallback.WithGroup(name),
		consoleLevels: levelRouter.consoleLevels,
	}
	if levelRouter.console != nil {
		newLevelRouter.console = levelRouter.console.WithGroup(name)
	}
	return newLevelRouter
}
