package exceptions

import (
	"net/http"
)

type HttpError struct {
	Message string
	Status  int
}

func NewHttpError(message string, status int) *HttpError {
	return &HttpError{
		Message: message,
		Status:  status,
	}
}

func (err *HttpError) Error() string {
	return err.Message
}

var (
	ConflictSqlError    = NewHttpError("Conflict creating resource.", http.StatusConflict)
	NotFoundError       = NewHttpError("Resource not found.", http.StatusNotFound)
	InternalServerError = NewHttpError("Ops, an error occurred. Try again later.", http.StatusInternalServerError)
	UnauthorizedError   = NewHttpError("Unauthorized.", http.StatusUnauthorized)
	ForbiddenError      = NewHttpError("Action forbidden.", http.StatusForbidden)
	BadRequestError     = NewHttpError("Invalid request.", http.StatusBadRequest)
)
