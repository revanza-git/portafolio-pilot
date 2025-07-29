package errors

import (
	"fmt"
	"net/http"
)

type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
	Status  int    `json:"-"`
}

func (e *AppError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Common error constructors
func New(code, message string, status int) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Status:  status,
	}
}

func BadRequest(message string) *AppError {
	return &AppError{
		Code:    "BAD_REQUEST",
		Message: message,
		Status:  http.StatusBadRequest,
	}
}

func Unauthorized(message string) *AppError {
	return &AppError{
		Code:    "UNAUTHORIZED",
		Message: message,
		Status:  http.StatusUnauthorized,
	}
}

func Forbidden(message string) *AppError {
	return &AppError{
		Code:    "FORBIDDEN",
		Message: message,
		Status:  http.StatusForbidden,
	}
}

func NotFound(resource string) *AppError {
	return &AppError{
		Code:    "NOT_FOUND",
		Message: fmt.Sprintf("%s not found", resource),
		Status:  http.StatusNotFound,
	}
}

func Conflict(message string) *AppError {
	return &AppError{
		Code:    "CONFLICT",
		Message: message,
		Status:  http.StatusConflict,
	}
}

func Internal(message string) *AppError {
	return &AppError{
		Code:    "INTERNAL_ERROR",
		Message: message,
		Status:  http.StatusInternalServerError,
	}
}

func ValidationError(details any) *AppError {
	return &AppError{
		Code:    "VALIDATION_ERROR",
		Message: "Validation failed",
		Details: details,
		Status:  http.StatusBadRequest,
	}
}

func DatabaseError(err error) *AppError {
	return &AppError{
		Code:    "DATABASE_ERROR",
		Message: "Database operation failed",
		Details: err.Error(),
		Status:  http.StatusInternalServerError,
	}
}

func ExternalServiceError(service string, err error) *AppError {
	return &AppError{
		Code:    "EXTERNAL_SERVICE_ERROR",
		Message: fmt.Sprintf("External service %s failed", service),
		Details: err.Error(),
		Status:  http.StatusServiceUnavailable,
	}
}