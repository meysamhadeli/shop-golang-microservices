package problemDetails

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"reflect"
	"time"
)

// ProblemDetail error struct
type ProblemDetail struct {
	Status    int       `json:"status,omitempty"`
	Title     string    `json:"title,omitempty"`
	Detail    string    `json:"detail,omitempty"`
	Type      string    `json:"type,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}

var mappers = map[reflect.Type]func() *ProblemDetail{}

// Error  Error() interface method
func (p *ProblemDetail) Error() string {
	return fmt.Sprintf("Error Title: %s - Error Status: %d - Error Detail: %s", p.Title, p.Status, p.Detail)
}

// WriteTo writes the JSON Problem to an HTTP Response Writer
func (p *ProblemDetail) WriteTo(w http.ResponseWriter) (int, error) {
	log.Error(p.Error())

	p.writeHeaderTo(w)
	return w.Write(p.json())
}

// Map map error to problem detail
func Map(err error, funcProb func() *ProblemDetail) {

	typeError := reflect.TypeOf(err).Elem()

	mappers[typeError] = funcProb
}

// ResolveProblemDetails retrieve error with format problem detail
func ResolveProblemDetails(err error) *ProblemDetail {

	typeError := reflect.TypeOf(err).Elem()

	problem := mappers[typeError]

	if problem != nil {
		return problem()
	}

	return &ProblemDetail{
		Type:      "https://httpstatuses.io/500",
		Status:    http.StatusInternalServerError,
		Detail:    err.Error(),
		Title:     "Internal Server Error",
		Timestamp: time.Now(),
	}
}

func (p *ProblemDetail) writeHeaderTo(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/problem+json")
	status := p.Status
	if status == 0 {
		status = http.StatusInternalServerError
	}

	w.WriteHeader(status)
}

func (p *ProblemDetail) json() []byte {
	b, _ := json.Marshal(&p)
	return b
}
