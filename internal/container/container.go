package container

import (
	"crud-task/internal/handlers"
	"crud-task/internal/repository"
	"crud-task/internal/services"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Container struct {
	TaskHandler *handlers.TaskHandler
}

func New(dbPool *pgxpool.Pool) *Container {
	baseRepository := repository.NewBaseRepository(dbPool)

	taskRepository := repository.NewTaskRepository(baseRepository)
	taskService := services.NewTaskService(taskRepository)
	taskHandler := handlers.NewTaskHandler(taskService)

	return &Container{
		TaskHandler: taskHandler,
	}
}
