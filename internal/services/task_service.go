package services

import (
	"context"
	"crud-task/internal/dto"
	"crud-task/internal/models"
	"crud-task/internal/repository"
	"github.com/jackc/pgx/v5"
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
	tx, err := taskService.taskRepository.DBPool.BeginTx(context, pgx.TxOptions{})
	if err != nil {
		return models.Task{}, err
	}
	defer finishTx(context, tx, &err)

	if task.Status == 0 {
		task.Status = models.Todo
	}

	created, err = taskService.taskRepository.Create(context, tx, task)
	return created, err
}

func (taskService *TaskService) Update(context context.Context, task models.Task) (err error) {
	tx, err := taskService.taskRepository.DBPool.BeginTx(context, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer finishTx(context, tx, &err)

	return taskService.taskRepository.Update(context, tx, task)
}

func (taskService *TaskService) Delete(context context.Context, id int64) (err error) {
	tx, err := taskService.taskRepository.DBPool.BeginTx(context, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer finishTx(context, tx, &err)

	return taskService.taskRepository.Delete(context, tx, id)
}

func finishTx(ctx context.Context, tx pgx.Tx, err *error) {
	if p := recover(); p != nil {
		_ = tx.Rollback(ctx)
		panic(p)
	} else if *err != nil {
		_ = tx.Rollback(ctx)
	} else {
		*err = tx.Commit(ctx)
	}
}

func (taskService *TaskService) GetStatusOptions() []dto.StatusOption {
	return []dto.StatusOption{
		{Value: int(models.Todo), Label: models.Todo.String()},
		{Value: int(models.InProgress), Label: models.InProgress.String()},
		{Value: int(models.Done), Label: models.Done.String()},
	}
}
