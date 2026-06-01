package errors

import (
	"errors"
	"fmt"
	"net/http"
)

type AppError struct {
	Message string
	Status  int
	Code    string
	Cause   error
}

func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Cause
}

func New(message string, status int, code string) *AppError {
	return &AppError{Message: message, Status: status, Code: code}
}

func Validation(message string) *AppError {
	return New(message, http.StatusUnprocessableEntity, "VALIDATION")
}

func Unauthorized(message string) *AppError {
	if message == "" {
		message = "Unauthorized"
	}
	return New(message, http.StatusUnauthorized, "UNAUTHORIZED")
}

func Forbidden(message string) *AppError {
	if message == "" {
		message = "Forbidden"
	}
	return New(message, http.StatusForbidden, "FORBIDDEN")
}

func TenantSuspended(message string) *AppError {
	if message == "" {
		message = "Tenant is suspended"
	}
	return New(message, http.StatusForbidden, "TENANT_SUSPENDED")
}

func NotFound(message string) *AppError {
	if message == "" {
		message = "Resource not found"
	}
	return New(message, http.StatusNotFound, "NOT_FOUND")
}

func Conflict(message string) *AppError {
	return New(message, http.StatusConflict, "CONFLICT")
}

func RateLimited(message string) *AppError {
	if message == "" {
		message = "Rate limit exceeded"
	}
	return New(message, http.StatusTooManyRequests, "RATE_LIMITED")
}

func Internal(message string, cause error) *AppError {
	return &AppError{
		Message: message,
		Status:  http.StatusInternalServerError,
		Code:    "INTERNAL",
		Cause:   cause,
	}
}

func AsAppError(err error) *AppError {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}
	return Internal("Internal server error", err)
}
