package appError

import (
	"fmt"
)

const (
	AppErrNotFound       = "NOT_FOUND"
	AppErrConflictRecord = "CONFLICT"
	AppErrFkViolation    = "FK_VIOLATION"
	AppErrFieldNotNull   = "NOT_NULL"
	AppErrCheckViolation = "CHECK_VIOLATION"
	AppErrInternal       = "INTERNAL"
)

type AppError struct {
	Status  int
	Code    string
	Message string
	Err     error
}

func (appErr *AppError) Error() string {
	if appErr.Err != nil {
		return fmt.Sprintf("%s: %v", appErr.Message, appErr.Err)
	}
	return appErr.Message
}

func New(status int, code string, message string, error error) *AppError {
	return &AppError{Status: status, Code: code, Message: message, Err: error}
}
