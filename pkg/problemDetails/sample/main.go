package main

import (
	"github.com/labstack/echo/v4"
	"github.com/meysamhadeli/shop-golang-microservices/pkg/problemDetails"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

func main() {
	e := echo.New()

	e.HTTPErrorHandler = ProblemDetailsHandler

	e.GET("/example1", example1)
	e.GET("/example2", example2)
	e.Logger.Fatal(e.Start(":5000"))
}

// example with built in problem details function error
func example1(c echo.Context) error {
	err := errors.New("We have a bad request in our endpoint")

	return problemDetails.BadRequestErr(err)
}

// example with create custom problem details error
func example2(c echo.Context) error {
	err := errors.New("We have a request timeout in our endpoint")

	return problemDetails.NewError(http.StatusRequestTimeout, err)
}

// middleware for handle problem details error on top of echo or gin or ...
func ProblemDetailsHandler(error error, c echo.Context) {

	// handle problem details with custom function error (it's optional)
	problemDetails.Map(http.StatusBadRequest, func() *problemDetails.ProblemDetail {
		return &problemDetails.ProblemDetail{
			Type:      "https://httpstatuses.io/400",
			Detail:    error.Error(),
			Title:     "bad-request",
			Timestamp: time.Now(),
		}
	})

	// handle problem details with custom function error (it's optional)
	problemDetails.Map(http.StatusRequestTimeout, func() *problemDetails.ProblemDetail {
		return &problemDetails.ProblemDetail{
			Type:      "https://httpstatuses.io/408",
			Status:    http.StatusRequestTimeout,
			Detail:    error.Error(),
			Title:     "request-timeout",
			Timestamp: time.Now(),
		}
	})

	// resolve problem details error from response in echo or gin or ...
	if !c.Response().Committed {
		if _, err := problemDetails.ResolveProblemDetails(c.Response(), error); err != nil {
			c.Logger().Error(err)
		}
	}
}
