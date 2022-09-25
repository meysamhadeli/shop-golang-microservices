package problemDetails

import (
	"encoding/json"
	"fmt"
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

// ProblemDetailError represents an error that occurred while handling a request.
type ProblemDetailError struct {
	Code     int         `json:"-"`
	Internal error       `json:"-"`
	Message  interface{} `json:"-"`
}

var mappers = map[int]func() *ProblemDetail{}

// WriteTo writes the JSON Problem to an HTTP Response Writer
func (p *ProblemDetail) writeTo(w http.ResponseWriter) (int, error) {
	p.writeHeaderTo(w)
	return w.Write(p.json())
}

// Map map error to problem detail
func Map(statusCode int, funcProblem func() *ProblemDetail) {
	mappers[statusCode] = funcProblem
}

// ResolveProblemDetails retrieve error with format problem detail
func ResolveProblemDetails(w http.ResponseWriter, err error) (int, error) {

	statusCode := err.(*ProblemDetailError).Code

	fmt.Println(mappers)

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

// Error makes it compatible with `error` interface.
func (he *ProblemDetailError) Error() string {
	return he.Internal.Error()
}

func NewError(code int, error error) *ProblemDetailError {
	newError := &ProblemDetailError{Code: code, Message: http.StatusText(code), Internal: error}
	return newError
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
