package problemDetails

import (
	"net/http"
	"time"
)

const (
	ErrBadRequestTitle          = "Bad Request"
	ErrConflictTitle            = "Conflict Error"
	ErrNotFoundTitle            = "Not Found"
	ErrUnauthorizedTitle        = "Unauthorized"
	ErrForbiddenTitle           = "Forbidden"
	ErrInternalServerErrorTitle = "Internal Server Error"
	ErrDomainTitle              = "Domain Model Error"
	ErrApplicationTitle         = "Application Service Error"
	ErrApiTitle                 = "Api Error"
	ErrRequestTimeoutTitle      = "Request Timeout"
)

func NewValidationProblemDetail(detail string, stackTrace string) ProblemDetailErr {
	validationError :=
		&problemDetail{
			Title:      ErrBadRequestTitle,
			Detail:     detail,
			Status:     http.StatusBadRequest,
			Type:       getDefaultType(http.StatusBadRequest),
			Timestamp:  time.Now(),
			StackTrace: stackTrace,
		}

	return validationError
}

func NewConflictProblemDetail(detail string, stackTrace string) ProblemDetailErr {
	return &problemDetail{
		Title:      ErrConflictTitle,
		Detail:     detail,
		Status:     http.StatusConflict,
		Type:       getDefaultType(http.StatusConflict),
		Timestamp:  time.Now(),
		StackTrace: stackTrace,
	}
}

func NewBadRequestProblemDetail(detail string, stackTrace string) ProblemDetailErr {
	return &problemDetail{
		Title:      ErrBadRequestTitle,
		Detail:     detail,
		Status:     http.StatusBadRequest,
		Type:       getDefaultType(http.StatusBadRequest),
		Timestamp:  time.Now(),
		StackTrace: stackTrace,
	}
}

func NewNotFoundErrorProblemDetail(detail string, stackTrace string) ProblemDetailErr {
	return &problemDetail{
		Title:      ErrNotFoundTitle,
		Detail:     detail,
		Status:     http.StatusNotFound,
		Type:       getDefaultType(http.StatusNotFound),
		Timestamp:  time.Now(),
		StackTrace: stackTrace,
	}
}

func NewUnAuthorizedErrorProblemDetail(detail string, stackTrace string) ProblemDetailErr {
	return &problemDetail{
		Title:      ErrUnauthorizedTitle,
		Detail:     detail,
		Status:     http.StatusUnauthorized,
		Type:       getDefaultType(http.StatusUnauthorized),
		Timestamp:  time.Now(),
		StackTrace: stackTrace,
	}
}

func NewForbiddenProblemDetail(detail string, stackTrace string) ProblemDetailErr {
	return &problemDetail{
		Title:      ErrForbiddenTitle,
		Detail:     detail,
		Status:     http.StatusForbidden,
		Type:       getDefaultType(http.StatusForbidden),
		Timestamp:  time.Now(),
		StackTrace: stackTrace,
	}
}

func NewInternalServerProblemDetail(detail string, stackTrace string) ProblemDetailErr {
	return &problemDetail{
		Title:      ErrInternalServerErrorTitle,
		Detail:     detail,
		Status:     http.StatusInternalServerError,
		Type:       getDefaultType(http.StatusInternalServerError),
		Timestamp:  time.Now(),
		StackTrace: stackTrace,
	}
}

func NewDomainProblemDetail(status int, detail string, stackTrace string) ProblemDetailErr {
	return &problemDetail{
		Title:      ErrDomainTitle,
		Detail:     detail,
		Status:     status,
		Type:       getDefaultType(http.StatusBadRequest),
		Timestamp:  time.Now(),
		StackTrace: stackTrace,
	}
}

func NewApplicationProblemDetail(status int, detail string, stackTrace string) ProblemDetailErr {
	return &problemDetail{
		Title:      ErrApplicationTitle,
		Detail:     detail,
		Status:     status,
		Type:       getDefaultType(status),
		Timestamp:  time.Now(),
		StackTrace: stackTrace,
	}
}

func NewApiProblemDetail(status int, detail string, stackTrace string) ProblemDetailErr {
	return &problemDetail{
		Title:      ErrApiTitle,
		Detail:     detail,
		Status:     status,
		Type:       getDefaultType(status),
		Timestamp:  time.Now(),
		StackTrace: stackTrace,
	}
}
