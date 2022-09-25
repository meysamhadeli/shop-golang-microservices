package problemDetails

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"reflect"
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

var mappers = map[reflect.Type]func() *ProblemDetail{}

// WriteTo writes the JSON Problem to an HTTP Response Writer
func (p *ProblemDetail) writeTo(w http.ResponseWriter) (int, error) {
	p.writeHeaderTo(w)
	return w.Write(p.json())
}

// Map map error to problem detail
func Map(error error, funcProblem func() *ProblemDetail) {

	typeError := reflect.TypeOf(error).Elem()

	mappers[typeError] = funcProblem
}

// ResolveEcho retrieve error with format problem detail
func ResolveEcho(res *echo.Response, err error) (int, error) {

	typeError := reflect.TypeOf(err).Elem()
	problem := mappers[typeError]

	if problem != nil {
		problem := problem()

		validationProblems(problem, err)

		val, err := problem.writeTo(res)

		if err != nil {
			return 0, err
		}

		return val, nil
	}

	defaultProblem := ProblemDetail{
		Type:      getDefaultType(http.StatusInternalServerError),
		Status:    http.StatusInternalServerError,
		Detail:    err.Error(),
		Timestamp: time.Now(),
	}

	val, err := defaultProblem.writeTo(res)

	if err != nil {
		return 0, err
	}

	return val, nil
}

func validationProblems(problem *ProblemDetail, err error) {
	if problem.Status == 0 {
		problem.Status = http.StatusInternalServerError
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
