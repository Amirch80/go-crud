package repository

import (
	"crud-task/pkg/appError"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"net/http"
)

const (
	codeUniqueViolation     = "23505"
	codeForeignKeyViolation = "23503"
	codeNotNullViolation    = "23502"
	codeCheckViolation      = "23514"
)

func MapError(error error) error {
	if error == nil {
		return nil
	}

	var appErr *appError.AppError
	if errors.As(error, &appErr) {
		return appErr
	}

	if errors.Is(error, pgx.ErrNoRows) {
		return appError.New(http.StatusNotFound, appError.AppErrNotFound, "record not found", error)
	}

	var postgresError *pgconn.PgError
	if errors.As(error, &postgresError) {
		switch postgresError.Code {
		case codeUniqueViolation:
			return appError.New(http.StatusConflict, appError.AppErrConflictRecord, "a record with these values already exists", error)
		case codeForeignKeyViolation:
			return appError.New(http.StatusBadRequest, appError.AppErrFkViolation, "referenced record does not exist", error)
		case codeNotNullViolation:
			return appError.New(http.StatusBadRequest, appError.AppErrFieldNotNull, "a required field is missing", error)
		case codeCheckViolation:
			return appError.New(http.StatusBadRequest, appError.AppErrCheckViolation, "a value violates a constraint", error)
		}
	}

	return appError.New(http.StatusInternalServerError, appError.AppErrInternal, "internal server error", error)
}
