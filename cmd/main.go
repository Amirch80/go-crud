package main

import (
	"context"
	"crud-task/internal/container"
	"crud-task/internal/database"
	"crud-task/internal/routes"
	"crud-task/pkg/logger"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	closeLogs, err := logger.Init()
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}

	exitCode := 0
	defer func() {
		closeLogs()
		os.Exit(exitCode)
	}()

	if err := run(); err != nil {
		logger.CustomLogger.Error("application error", "error", err.Error())
		exitCode = 1
		return
	}
}

func run() error {
	dbPool, err := database.Connect()
	if err != nil {
		return err
	}
	if dbPool == nil {
		return errors.New("dbPool is nil — Connect() returned nothing")
	}
	defer database.Close()

	appContainer := container.New(dbPool)
	server := &http.Server{
		Addr:    ":8080",
		Handler: routes.Routes(appContainer),
	}

	serverErr := make(chan error, 1)
	go func() {
		logger.CustomLogger.Info("server started", slog.String("address", server.Addr))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErr:
		return err
	case <-quit:
		logger.CustomLogger.Info("shutting down server...")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		return err
	}

	logger.CustomLogger.Info("server exited properly")
	return nil
}
