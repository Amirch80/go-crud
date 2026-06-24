package main

import (
	"context"
	"crud-task/internal/container"
	"crud-task/internal/database"
	"crud-task/internal/routes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	dbPool := database.Connect()

	if dbPool == nil {
		log.Fatal("dbPool is nil — Connect() returned nothing")
	}
	defer database.Close()

	appContainer := container.New(dbPool)
	handler := routes.Routes(appContainer)

	server := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	go func() {
		fmt.Println("Server started on :8080")
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	fmt.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("Shutdown error: %s\n", err)
	}

	fmt.Println("Server exited properly")
}
