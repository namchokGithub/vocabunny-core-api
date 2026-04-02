package helper

import (
	"errors"
	"net/http"
)

type AppError struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	StatusCode int    `json:"-"`
	Err        error  `json:"-"`
}

func (e *AppError) Error() string {
	if e == nil {
		return ""
	}

	return e.Message
}

func (e *AppError) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Err
}

func NewAppError(statusCode int, code, message string, err error) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
		Err:        err,
	}
}

func BadRequest(code, message string, err error) *AppError {
	return NewAppError(http.StatusBadRequest, code, message, err)
}

func NotFound(code, message string, err error) *AppError {
	return NewAppError(http.StatusNotFound, code, message, err)
}

func Conflict(code, message string, err error) *AppError {
	return NewAppError(http.StatusConflict, code, message, err)
}

func Internal(code, message string, err error) *AppError {
	return NewAppError(http.StatusInternalServerError, code, message, err)
}

func Unauthorized(code, message string, err error) *AppError {
	return NewAppError(http.StatusUnauthorized, code, message, err)
}

func Forbidden(code, message string, err error) *AppError {
	return NewAppError(http.StatusForbidden, code, message, err)
}

func AsAppError(err error) *AppError {
	if err == nil {
		return nil
	}

	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}

	return Internal("internal_error", "internal server error", err)
}
