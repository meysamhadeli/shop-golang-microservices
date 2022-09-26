package problemDetails

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

// ProblemDetail error struct
type ProblemDetail struct {
	Status     int       `json:"status,omitempty"`
	Title      string    `json:"title,omitempty"`
	Detail     string    `json:"detail,omitempty"`
	Type       string    `json:"type,omitempty"`
	Timestamp  time.Time `json:"timestamp,omitempty"`
	StackTrace string    `json:"stackTrace,omitempty"`
}

// ProblemDetailError represents an error in problem details
type problemDetailError struct {
	Code    int   `json:"-"`
	Details error `json:"-"`
}

var mappers = map[int]func() *ProblemDetail{}

// WriteTo writes the JSON Problem to an HTTP Response Writer
func (p *ProblemDetail) writeTo(w http.ResponseWriter) (int, error) {
	p.writeHeaderTo(w)
	return w.Write(p.json())
}

// Map map error to problem details error
func Map(statusCode int, funcProblem func() *ProblemDetail) {
	mappers[statusCode] = funcProblem
}

// ResolveProblemDetails retrieve and resolve error with format problem details error
func ResolveProblemDetails(w http.ResponseWriter, err error) (int, error) {

	var statusCode int

	var pe *problemDetailError

	if errors.As(err, &pe) == false {
		statusCode = http.StatusInternalServerError
	} else {
		statusCode = err.(*problemDetailError).Code
	}

	problem := mappers[statusCode]

	if problem != nil {
		problem := problem()

		validationProblems(problem, err, statusCode)

		val, err := problem.writeTo(w)

		if err != nil {
			return 0, err
		}

		return val, err
	}

	defaultProblem := func() *ProblemDetail {
		return &ProblemDetail{
			Type:      getDefaultType(statusCode),
			Status:    statusCode,
			Detail:    err.Error(),
			Timestamp: time.Now(),
			Title:     http.StatusText(statusCode),
		}
	}

	val, err := defaultProblem().writeTo(w)

	if err != nil {
		return 0, err
	}

	return val, nil
}

// Error makes error compatible with `error` interface.
func (p *problemDetailError) Error() string {
	if p.Details == nil {
		return fmt.Sprintf("code=%d", p.Code)
	}
	return p.Details.Error()
}

// NewError make custom error compatible with `error` interface.
func NewError(code int, error error) *problemDetailError {
	newError := &problemDetailError{Code: code, Details: error}
	return newError
}

// BadRequestErr make badRequest error compatible with `error` interface.
func BadRequestErr(error error) *problemDetailError {
	return NewError(http.StatusBadRequest, error)
}

// InternalServerErr make internalServer error compatible with `error` interface.
func InternalServerErr(error error) *problemDetailError {
	return NewError(http.StatusInternalServerError, error)
}

// NotFoundErr make notFound error compatible with `error` interface.
func NotFoundErr(error error) *problemDetailError {
	return NewError(http.StatusNotFound, error)
}

// UnauthorizedErr make unauthorized error compatible with `error` interface.
func UnauthorizedErr(error error) *problemDetailError {
	return NewError(http.StatusUnauthorized, error)
}

// ForbiddenErr make forbidden error compatible with `error` interface.
func ForbiddenErr(error error) *problemDetailError {
	return NewError(http.StatusForbidden, error)
}

// UnsupportedMediaTypeErr make unsupportedMediaType error compatible with `error` interface.
func UnsupportedMediaTypeErr(error error) *problemDetailError {
	return NewError(http.StatusUnsupportedMediaType, error)
}

// BadGatewayErr make badGateway error compatible with `error` interface.
func BadGatewayErr(error error) *problemDetailError {
	return NewError(http.StatusBadGateway, error)
}

func validationProblems(problem *ProblemDetail, err error, statusCode int) {

	if problem.Status == 0 {
		problem.Status = statusCode
	}
	if problem.Timestamp.IsZero() {
		problem.Timestamp = time.Now()
	}
	if problem.Detail == "" {
		problem.Detail = err.Error()
	}
	if problem.Type == "" {
		problem.Type = getDefaultType(problem.Status)
	}
	if problem.Title == "" {
		problem.Title = http.StatusText(problem.Status)
	}
}

func (p *ProblemDetail) writeHeaderTo(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/problem+json")

	w.WriteHeader(p.Status)
}

func (p *ProblemDetail) json() []byte {
	res, _ := json.Marshal(&p)
	return res
}

func getDefaultType(statusCode int) string {
	return fmt.Sprintf("https://httpstatuses.io/%d", statusCode)
}
