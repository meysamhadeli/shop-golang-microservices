package customErrors

import (
	"github.com/pkg/errors"
)

func NewMarshalingError(message string) error {
	internal := NewInternalServerError(message)
	customErr := GetCustomError(internal)
	ue := &marshalingError{
		InternalServerError: customErr.(InternalServerError),
	}
	stackErr := errors.WithStack(ue)

	return stackErr
}

func NewMarshalingErrorWrap(err error, message string) error {
	internal := NewInternalServerErrorWrap(err, message)
	customErr := GetCustomError(internal)
	ue := &marshalingError{
		InternalServerError: customErr.(InternalServerError),
	}
	stackErr := errors.WithStack(ue)

	return stackErr
}

type marshalingError struct {
	InternalServerError
}

type MarshalingError interface {
	InternalServerError
	IsMarshalingError() bool
}

func (m *marshalingError) IsMarshalingError() bool {
	return true
}

func IsMarshalingError(err error) bool {
	var me MarshalingError

	//us, ok := grpc_errors.Cause(err).(MarshalingError)
	if errors.As(err, &me) {
		return me.IsMarshalingError()
	}

	return false
}
