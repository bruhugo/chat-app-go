package exceptions

import (
	"errors"
	"net/http"
)

type HttpError struct {
	Err     error
	Message string
	Status  int
}

func NewHttpError(err error, status int) *HttpError {
	return &HttpError{
		Err:    err,
		Status: status,
	}
}

func NewHttpErrorWithMessage(err error, status int, message string) *HttpError {
	return &HttpError{
		Err:     err,
		Message: message,
		Status:  status,
	}
}

func (err *HttpError) Error() string {
	return err.Err.Error()
}

var (
	ConflictSqlError    = NewHttpError(errors.New("Conflict creating resource."), http.StatusConflict)
	NotFoundError       = NewHttpError(errors.New("Resource not found."), http.StatusNotFound)
	InternalServerError = NewHttpError(errors.New("Ops, an error occurred. Try again later."), http.StatusNotFound)
)
