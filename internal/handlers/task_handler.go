package handlers

import (
	"crud-task/internal/models"
	"crud-task/internal/services"
	"crud-task/pkg"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5"
	"net/http"
	"strconv"
)

type TaskHandler struct {
	taskService *services.TaskService
}

func NewTaskHandler(taskService *services.TaskService) *TaskHandler {
	return &TaskHandler{
		taskService: taskService,
	}
}

func (taskHandler *TaskHandler) All(writer http.ResponseWriter, request *http.Request) {
	tasks, err := taskHandler.taskService.All(request.Context())
	if err != nil {
		pkg.ResponseError(writer, http.StatusInternalServerError, err.Error())
		return
	}
	pkg.ResponseJson(writer, http.StatusOK, tasks)
}

func (taskHandler *TaskHandler) Show(writer http.ResponseWriter, request *http.Request) {
	idStr := request.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		pkg.ResponseError(writer, http.StatusNotFound, "Task not found")
		return
	}

	task, err := taskHandler.taskService.Show(request.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			pkg.ResponseError(writer, http.StatusNotFound, "Task not found")
			return
		}
		pkg.ResponseError(writer, http.StatusInternalServerError, err.Error())
		return
	}
	pkg.ResponseJson(writer, http.StatusOK, task)
}

func (taskHandler *TaskHandler) Create(writer http.ResponseWriter, request *http.Request) {
	var task models.Task
	if err := json.NewDecoder(request.Body).Decode(&task); err != nil {
		pkg.ResponseError(writer, http.StatusBadRequest, "invalid body")
		return
	}

	created, err := taskHandler.taskService.Create(request.Context(), task)
	if err != nil {
		pkg.ResponseError(writer, http.StatusInternalServerError, err.Error())
		return
	}
	pkg.ResponseJson(writer, http.StatusCreated, created)
}

func (taskHandler *TaskHandler) Update(writer http.ResponseWriter, request *http.Request) {
	idStr := request.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		pkg.ResponseError(writer, http.StatusBadRequest, "Invalid id")
		return
	}

	var task models.Task
	if err := json.NewDecoder(request.Body).Decode(&task); err != nil {
		pkg.ResponseError(writer, http.StatusBadRequest, "Invalid body")
		return
	}
	task.Id = id

	if err := taskHandler.taskService.Update(request.Context(), task); err != nil {
		pkg.ResponseError(writer, http.StatusInternalServerError, err.Error())
		return
	}
	pkg.ResponseJson(writer, http.StatusOK, task)
}

func (taskHandler *TaskHandler) Delete(writer http.ResponseWriter, request *http.Request) {
	idStr := request.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		pkg.ResponseError(writer, http.StatusBadRequest, "Invalid id")
		return
	}

	if err := taskHandler.taskService.Delete(request.Context(), id); err != nil {
		pkg.ResponseError(writer, http.StatusInternalServerError, err.Error())
		return
	}
	pkg.ResponseJson(writer, http.StatusNoContent, nil)
}

func (taskHandler *TaskHandler) StatusOptions(writer http.ResponseWriter, _ *http.Request) {
	options := taskHandler.taskService.GetStatusOptions()
	pkg.ResponseJson(writer, http.StatusOK, options)
}
