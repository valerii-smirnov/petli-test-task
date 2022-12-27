package messages

import "net/http"

type InternalServerError struct {
	Code    int    `json:"code" example:"500"`
	Message string `json:"message" example:"something went wrong"`
	Error   string `json:"error" example:"error message"`
}

func NewInternalServerError(message string, error error) InternalServerError {
	return InternalServerError{
		Code:    http.StatusInternalServerError,
		Message: message,
		Error:   error.Error(),
	}
}

type BadRequestError struct {
	Code    int               `json:"code" example:"400"`
	Message string            `json:"message" example:"validation error"`
	Errors  map[string]string `json:"errors,omitempty"`
}

func NewBadRequestError(message string, errors map[string]string) BadRequestError {
	return BadRequestError{
		Code:    http.StatusBadRequest,
		Message: message,
		Errors:  errors,
	}
}

type NotFoundError struct {
	Code    int    `json:"code" example:"404"`
	Message string `json:"message" example:"dog not found"`
}

func NewNotFoundError(message string) NotFoundError {
	return NotFoundError{
		Code:    http.StatusNotFound,
		Message: message,
	}
}

type UnauthenticatedError struct {
	Code    int    `json:"code" example:"401"`
	Message string `json:"message" example:"unauthorized"`
}

func NewUnauthenticatedError(message string) UnauthenticatedError {
	return UnauthenticatedError{
		Code:    http.StatusUnauthorized,
		Message: message,
	}
}

type ForbiddenError struct {
	Code    int    `json:"code" example:"403"`
	Message string `json:"message" example:"unauthorized"`
}

func NewForbiddenError(message string) ForbiddenError {
	return ForbiddenError{
		Code:    http.StatusForbidden,
		Message: message,
	}
}

type ConflictError struct {
	Code    int    `json:"code" example:"409"`
	Message string `json:"message" example:"dog already exists"`
}

func NewConflictError(message string) *ConflictError {
	return &ConflictError{
		Code:    http.StatusConflict,
		Message: message,
	}
}
