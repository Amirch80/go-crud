package handlers

import (
	"crud-task/internal/dto"
	"crud-task/internal/models"
	"crud-task/internal/services"
	"crud-task/pkg"
	"encoding/json"
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
	response := pkg.Response{}

	tasks, err := taskHandler.taskService.All(request.Context())
	if err != nil {
		pkg.ResponseFromError(writer, err)
		return
	}
	response = pkg.Response{
		Data: map[string][]models.Task{
			"tasks": tasks,
		},
	}
	pkg.ResponseJson(writer, http.StatusOK, &response)
}

func (taskHandler *TaskHandler) Show(writer http.ResponseWriter, request *http.Request) {
	idStr := request.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		pkg.ResponseError(writer, http.StatusNotFound, "task not found")
		return
	}

	task, err := taskHandler.taskService.Show(request.Context(), id)
	if err != nil {
		pkg.ResponseFromError(writer, err)
		return
	}
	response := pkg.Response{
		Data: map[string]models.Task{
			"task": task,
		},
	}
	pkg.ResponseJson(writer, http.StatusOK, &response)
}

func (taskHandler *TaskHandler) Create(writer http.ResponseWriter, request *http.Request) {
	var task models.Task
	if err := json.NewDecoder(request.Body).Decode(&task); err != nil {
		pkg.ResponseError(writer, http.StatusBadRequest, "invalid body")
		return
	}

	validationErrors, err := pkg.Validate(task)
	if err != nil {
		pkg.ResponseError(writer, http.StatusInternalServerError, "internal server error")
		return
	}

	if validationErrors != nil {
		pkg.ResponseValidationError(writer, http.StatusUnprocessableEntity, validationErrors)
		return
	}

	created, err := taskHandler.taskService.Create(request.Context(), task)
	if err != nil {
		pkg.ResponseFromError(writer, err)
		return
	}
	response := pkg.Response{
		Data: map[string]models.Task{
			"task": created,
		},
	}
	pkg.ResponseJson(writer, http.StatusCreated, &response)
}

func (taskHandler *TaskHandler) Update(writer http.ResponseWriter, request *http.Request) {
	idStr := request.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		pkg.ResponseError(writer, http.StatusBadRequest, "invalid id")
		return
	}

	var task models.Task
	if err := json.NewDecoder(request.Body).Decode(&task); err != nil {
		pkg.ResponseError(writer, http.StatusBadRequest, "invalid body")
		return
	}
	task.Id = id

	validationErrors, err := pkg.Validate(task)
	if err != nil {
		pkg.ResponseError(writer, http.StatusInternalServerError, "internal server error")
		return
	}

	if validationErrors != nil {
		pkg.ResponseValidationError(writer, http.StatusUnprocessableEntity, validationErrors)
		return
	}

	if err := taskHandler.taskService.Update(request.Context(), task); err != nil {
		pkg.ResponseFromError(writer, err)
		return
	}
	response := pkg.Response{
		Data: map[string]models.Task{
			"task": task,
		},
	}
	pkg.ResponseJson(writer, http.StatusOK, &response)
}

func (taskHandler *TaskHandler) Delete(writer http.ResponseWriter, request *http.Request) {
	idStr := request.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		pkg.ResponseError(writer, http.StatusBadRequest, "invalid id")
		return
	}

	if err := taskHandler.taskService.Delete(request.Context(), id); err != nil {
		pkg.ResponseFromError(writer, err)
		return
	}
	pkg.ResponseJson(writer, http.StatusNoContent, nil)
}

func (taskHandler *TaskHandler) StatusOptions(writer http.ResponseWriter, _ *http.Request) {
	options := taskHandler.taskService.GetStatusOptions()
	response := pkg.Response{
		Data: map[string][]dto.StatusOption{
			"statuses": options,
		},
	}
	pkg.ResponseJson(writer, http.StatusOK, &response)
}
