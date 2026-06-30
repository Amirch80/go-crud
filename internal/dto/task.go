package dto

import (
	"crud-task/internal/models"
	"net/url"
	"strconv"
	"strings"
)

type CreateTaskRequest struct {
	Title       string        `json:"title" validate:"required,min=5,max=255"`
	Description string        `json:"description" validate:"max=1000"`
	Status      models.Status `json:"status" validate:"required,oneof=1 2 3"`
}

type UpdateTaskRequest struct {
	Title       string        `json:"title" validate:"required,min=5,max=255"`
	Description string        `json:"description" validate:"max=1000"`
	Status      models.Status `json:"status" validate:"required,oneof=1 2 3"`
}

type ListTaskQuery struct {
	Page    int
	PerPage int
	Sort    string
	Status  *models.Status
}

func ConvertQueryToListTaskQuery(values url.Values) ListTaskQuery {
	listTaskQuery := ListTaskQuery{
		Page:    1,
		PerPage: 20,
		Sort:    "desc",
		Status:  nil,
	}

	if page, err := strconv.Atoi(values.Get("page")); err != nil && page > 0 {
		listTaskQuery.Page = page
	}

	if sort := strings.ToLower(values.Get("sort")); sort == "asc" || sort == "desc" {
		listTaskQuery.Sort = sort
	}

	if statusSlug := values.Get("status"); statusSlug != "" {
		status, ok := models.ParseStatusSlug(statusSlug)
		if !ok {
			return listTaskQuery
		}
		listTaskQuery.Status = &status
	}

	return listTaskQuery
}
