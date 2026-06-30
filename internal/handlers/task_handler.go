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

	listTaskQuery := dto.ConvertQueryToListTaskQuery(request.URL.Query())
	tasks, total, err := taskHandler.taskService.All(request.Context(), listTaskQuery)
	if err != nil {
		pkg.ResponseFromError(writer, err)
		return
	}
	response.Data = map[string]any{
		"tasks":    tasks,
		"total":    total,
		"per_page": listTaskQuery.PerPage,
		"page":     listTaskQuery.Page,
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
	var createTaskRequest dto.CreateTaskRequest
	if err := json.NewDecoder(request.Body).Decode(&createTaskRequest); err != nil {
		pkg.ResponseError(writer, http.StatusBadRequest, "invalid body")
		return
	}

	validationErrors, err := pkg.Validate(createTaskRequest)
	if err != nil {
		pkg.ResponseError(writer, http.StatusInternalServerError, "internal server error")
		return
	}

	if validationErrors != nil {
		pkg.ResponseValidationError(writer, http.StatusUnprocessableEntity, validationErrors)
		return
	}

	task := models.Task{
		Title:       createTaskRequest.Title,
		Description: createTaskRequest.Description,
		Status:      createTaskRequest.Status,
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

	var updateTaskRequest dto.UpdateTaskRequest

	if err := json.NewDecoder(request.Body).Decode(&updateTaskRequest); err != nil {
		pkg.ResponseError(writer, http.StatusBadRequest, "invalid body")
		return
	}

	validationErrors, err := pkg.Validate(updateTaskRequest)
	if err != nil {
		pkg.ResponseError(writer, http.StatusInternalServerError, "internal server error")
		return
	}

	if validationErrors != nil {
		pkg.ResponseValidationError(writer, http.StatusUnprocessableEntity, validationErrors)
		return
	}

	task := models.Task{
		Id:          id,
		Title:       updateTaskRequest.Title,
		Description: updateTaskRequest.Description,
		Status:      updateTaskRequest.Status,
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

func (taskHandler *TaskHandler) SoftDelete(writer http.ResponseWriter, request *http.Request) {
	idStr := request.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		pkg.ResponseError(writer, http.StatusBadRequest, "invalid id")
		return
	}

	if err := taskHandler.taskService.SoftDelete(request.Context(), id); err != nil {
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
