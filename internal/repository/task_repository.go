package repository

import (
	"context"
	"crud-task/internal/dto"
	"crud-task/internal/models"
	"crud-task/pkg/appError"
	"fmt"
	"github.com/doug-martin/goqu/v9"
	"net/http"
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

func (taskRepository *TaskRepository) All(context context.Context, listTaskQuery dto.ListTaskQuery) ([]models.Task, int64, error) {
	query := dialect.From("tasks")

	if listTaskQuery.Status != nil {
		query = query.Where(
			goqu.C("status").Eq(listTaskQuery.Status),
		)
	}

	countSql, countArgs, err := query.Select(goqu.COUNT("id").As("total")).Prepared(true).ToSQL()
	if err != nil {
		return nil, 0, fmt.Errorf("build count tasks query: %w", err)
	}
	var total int64
	if err := taskRepository.DBPool.QueryRow(context, countSql, countArgs...).Scan(&total); err != nil {
		return nil, 0, MapError(err)
	}

	if listTaskQuery.Sort == "desc" {
		query = query.Order(goqu.C("created_at").Desc())
	} else {
		query = query.Order(goqu.C("created_at").Asc())
	}

	offset := (listTaskQuery.Page - 1) * listTaskQuery.PerPage

	sql, args, err := query.Select("id", "title", "description", "status", "created_at", "updated_at").
		Where(goqu.C("deleted_at").IsNull()).
		Offset(uint(offset)).
		Limit(uint(listTaskQuery.PerPage)).
		Prepared(true).
		ToSQL()

	if err != nil {
		return nil, 0, fmt.Errorf("build all tasks query: %w", err)
	}

	rows, err := taskRepository.DBPool.Query(context, sql, args...)

	if err != nil {
		return nil, 0, MapError(err)
	}

	defer rows.Close()

	var tasks []models.Task

	for rows.Next() {
		var task models.Task
		if err := rows.Scan(
			&task.Id, &task.Title, &task.Description,
			&task.Status, &task.CreatedAt, &task.UpdatedAt,
		); err != nil {
			return nil, 0, MapError(err)
		}
		tasks = append(tasks, task)
	}

	return tasks, total, MapError(rows.Err())
}

func (taskRepository *TaskRepository) Show(ctx context.Context, id int64) (models.Task, error) {
	sql, args, err := dialect.
		From("tasks").
		Select("id", "title", "description", "status", "created_at", "updated_at").
		Where(
			goqu.C("id").Eq(id),
			goqu.C("deleted_at").IsNull(),
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

	return task, MapError(err)
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
		return models.Task{}, err
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
		Where(
			goqu.C("id").Eq(task.Id),
			goqu.C("deleted_at").IsNull(),
		).
		Prepared(true).
		ToSQL()

	if err != nil {
		return fmt.Errorf("build update task query: %w", err)
	}

	tag, err := query.Exec(context, sql, args...)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return appError.New(http.StatusNotFound, appError.AppErrNotFound, fmt.Sprintf("task with id %d not found", task.Id), nil)
	}
	return nil
}

// Delete TODO this action requires delete permission
func (taskRepository *TaskRepository) Delete(context context.Context, query DBTX, id int64) error {
	sql, args, err := dialect.
		Delete("tasks").
		Where(goqu.C("id").Eq(id)).
		Prepared(true).
		ToSQL()
	if err != nil {
		return fmt.Errorf("build delete task query: %w", err)
	}

	tag, err := query.Exec(context, sql, args...)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return appError.New(http.StatusNotFound, appError.AppErrNotFound, fmt.Sprintf("task with id %d not found", id), nil)
	}
	return nil
}

func (taskRepository *TaskRepository) SoftDelete(context context.Context, query DBTX, id int64) error {
	sql, args, err := dialect.
		Update("tasks").
		Where(
			goqu.C("id").Eq(id),
			goqu.C("deleted_at").IsNull(),
		).
		Set(goqu.Record{
			"deleted_at": time.Now(),
		}).
		Prepared(true).
		ToSQL()
	if err != nil {
		return fmt.Errorf("build delete task query: %w", err)
	}

	tag, err := query.Exec(context, sql, args...)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return appError.New(http.StatusNotFound, appError.AppErrNotFound, fmt.Sprintf("task with id %d not found", id), nil)
	}
	return nil
}
