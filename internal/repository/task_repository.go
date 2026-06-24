package repository

import (
	"context"
	"crud-task/internal/models"
	"fmt"
	"github.com/doug-martin/goqu/v9"
	"time"
)

type TaskRepository struct {
	*BaseRepository
}

func NewTaskRepository(repository *BaseRepository) *TaskRepository {
	return &TaskRepository{
		BaseRepository: repository,
	}
}

func (taskRepository *TaskRepository) All(context context.Context) ([]models.Task, error) {
	sql, args, err := dialect.From("tasks").Prepared(true).ToSQL()
	if err != nil {
		return nil, fmt.Errorf("build all tasks query: %w", err)
	}

	rows, err := taskRepository.DBPool.Query(context, sql, args...)

	if err != nil {
		return nil, fmt.Errorf("query all tasks: %w", err)
	}

	defer rows.Close()

	var tasks []models.Task

	for rows.Next() {
		var task models.Task
		if err := rows.Scan(
			&task.Id, &task.Title, &task.Description,
			&task.Status, &task.CreatedAt, &task.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan task: %w", err)
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate tasks: %w", err)
	}
	return tasks, nil
}

func (taskRepository *TaskRepository) Show(ctx context.Context, id int64) (models.Task, error) {
	sql, args, err := dialect.
		From("tasks").
		Where(
			goqu.C("id").Eq(id),
		).
		Prepared(true).
		ToSQL()
	if err != nil {
		return models.Task{}, fmt.Errorf("build show task query: %w", err)
	}

	var task models.Task
	err = taskRepository.DBPool.QueryRow(ctx, sql, args...).Scan(
		&task.Id, &task.Title, &task.Description,
		&task.Status, &task.CreatedAt, &task.UpdatedAt,
	)
	if err != nil {
		return models.Task{}, fmt.Errorf("show task: %w", err)
	}
	return task, nil
}

func (taskRepository *TaskRepository) Create(context context.Context, query DBTX, task models.Task) (models.Task, error) {
	sql, args, err := dialect.
		Insert("tasks").
		Rows(goqu.Record{
			"title":       task.Title,
			"description": task.Description,
			"status":      task.Status,
			"created_at":  goqu.L("NOW()"),
			"updated_at":  goqu.L("NOW()"),
		}).
		Returning("id", "created_at", "updated_at").
		Prepared(true).
		ToSQL()
	if err != nil {
		return models.Task{}, fmt.Errorf("build create task query: %w", err)
	}

	err = query.QueryRow(context, sql, args...).Scan(&task.Id, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		return models.Task{}, fmt.Errorf("create task: %w", err)
	}
	return task, nil
}

func (taskRepository *TaskRepository) Update(context context.Context, query DBTX, task models.Task) error {
	sql, args, err := dialect.
		Update("tasks").
		Set(goqu.Record{
			"title":       task.Title,
			"description": task.Description,
			"status":      task.Status,
			"updated_at":  time.Now(),
		}).
		Where(goqu.C("id").Eq(task.Id)).
		Prepared(true).
		ToSQL()

	if err != nil {
		return fmt.Errorf("build update task query: %w", err)
	}

	tag, err := query.Exec(context, sql, args...)
	if err != nil {
		return fmt.Errorf("update task: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("task with id %d not found", task.Id)
	}
	return nil
}

func (taskRepository *TaskRepository) Delete(context context.Context, q DBTX, id int64) error {
	sql, args, err := dialect.
		Delete("tasks").
		Where(goqu.C("id").Eq(id)).
		Prepared(true).
		ToSQL()
	if err != nil {
		return fmt.Errorf("build delete task query: %w", err)
	}

	tag, err := q.Exec(context, sql, args...)
	if err != nil {
		return fmt.Errorf("delete task: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("task with id %d not found", id)
	}
	return nil
}
