package services

import (
	"context"
	"crud-task/internal/dto"
	"crud-task/internal/models"
	"crud-task/internal/repository"
)

type TaskService struct {
	taskRepository *repository.TaskRepository
}

func NewTaskService(taskRepository *repository.TaskRepository) *TaskService {
	return &TaskService{
		taskRepository: taskRepository,
	}
}

func (taskService *TaskService) All(context context.Context) ([]models.Task, error) {
	return taskService.taskRepository.All(context)
}

func (taskService *TaskService) Show(context context.Context, id int64) (models.Task, error) {
	return taskService.taskRepository.Show(context, id)
}

func (taskService *TaskService) Create(context context.Context, task models.Task) (created models.Task, err error) {
	if task.Status == 0 {
		task.Status = models.Todo
	}
	err = taskService.taskRepository.RunInTx(context, func(query repository.DBTX) error {
		created, err = taskService.taskRepository.Create(context, query, task)
		return err
	})
	return created, err
}

func (taskService *TaskService) Update(context context.Context, task models.Task) (err error) {
	err = taskService.taskRepository.RunInTx(context, func(query repository.DBTX) error {
		return taskService.taskRepository.Update(context, query, task)
	})
	return err
}

func (taskService *TaskService) Delete(context context.Context, id int64) (err error) {
	err = taskService.taskRepository.RunInTx(context, func(query repository.DBTX) error {
		return taskService.taskRepository.Delete(context, query, id)
	})
	return err
}

func (taskService *TaskService) GetStatusOptions() []dto.StatusOption {
	return []dto.StatusOption{
		{Value: int(models.Todo), Label: models.Todo.String()},
		{Value: int(models.InProgress), Label: models.InProgress.String()},
		{Value: int(models.Done), Label: models.Done.String()},
	}
}
