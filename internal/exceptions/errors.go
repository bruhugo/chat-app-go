package exceptions

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
