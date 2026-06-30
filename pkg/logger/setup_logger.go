package logger

import (
	"log/slog"
	"os"
)

var CustomLogger *slog.Logger

func openFile(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
}

func setupLogger() (*slog.Logger, func(), error) {
	basePath := "./logs/"

	if os.Getenv("APP_ENV") == "development" {
		basePath = "../../logs/"
	}

	paths := map[slog.Level]string{
		slog.LevelDebug: basePath + "debug.log",
		slog.LevelInfo:  basePath + "info.log",
		slog.LevelWarn:  basePath + "warn.log",
		slog.LevelError: basePath + "error.log",
	}

	var opened []*os.File

	closeAll := func() {
		for _, file := range opened {
			file.Close()
		}
	}

	options := &slog.HandlerOptions{Level: slog.LevelDebug} //The lowest value should be set

	files := make(map[slog.Level]slog.Handler, len(paths))

	for level, path := range paths {
		file, err := openFile(path)
		if err != nil {
			closeAll()
			return nil, nil, err
		}
		opened = append(opened, file)
		files[level] = slog.NewJSONHandler(file, options)
	}

	other, err := openFile(basePath + "other.log")
	if err != nil {
		other.Close()
		return nil, nil, err
	}
	opened = append(opened, other)

	levelRouter := &LevelRouter{
		files:    files,
		fallback: slog.NewJSONHandler(other, options),
		console:  slog.NewTextHandler(os.Stdout, options),
		consoleLevels: map[slog.Level]bool{
			slog.LevelDebug: false,
			slog.LevelInfo:  false,
			slog.LevelWarn:  false,
			slog.LevelError: false,
		},
	}

	return slog.New(levelRouter), closeAll, nil
}

func Init() (func(), error) {
	logger, closeLogs, err := setupLogger()
	if err != nil {
		return nil, err
	}
	CustomLogger = logger
	slog.SetDefault(logger)
	return closeLogs, nil
}
