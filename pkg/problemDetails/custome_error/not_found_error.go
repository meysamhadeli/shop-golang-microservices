package customErrors

import (
	"github.com/pkg/errors"
	"net/http"
)

func NewNotFoundError(message string) error {
	ne := &notFoundError{
		CustomError: NewCustomError(nil, http.StatusNotFound, message),
	}
	stackErr := errors.WithStack(ne)

	return stackErr
}

func NewNotFoundErrorWrap(err error, message string) error {
	ne := &notFoundError{
		CustomError: NewCustomError(err, http.StatusNotFound, message),
	}
	stackErr := errors.WithStack(ne)

	return stackErr
}

type notFoundError struct {
	CustomError
}

type NotFoundError interface {
	CustomError
	IsNotFoundError() bool
}

func (n *notFoundError) IsNotFoundError() bool {
	return true
}

func IsNotFoundError(err error) bool {
	var notFoundError NotFoundError
	//us, ok := grpc_errors.Cause(err).(NotFoundError)
	if errors.As(err, &notFoundError) {
		return notFoundError.IsNotFoundError()
	}

	return false
}
