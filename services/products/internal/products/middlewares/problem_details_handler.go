package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/meysamhadeli/problem-details"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func ProblemDetailsHandler(error error, c echo.Context) {

	problem.Map(http.StatusBadRequest, func() *problem.ProblemDetail {
		return &problem.ProblemDetail{
			Type:      "https://httpstatuses.io/400",
			Detail:    error.Error(),
			Title:     "bad-request",
			Timestamp: time.Now(),
		}
	})

	problem.Map(http.StatusForbidden, func() *problem.ProblemDetail {
		return &problem.ProblemDetail{
			Type:      "https://httpstatuses.io/403",
			Status:    http.StatusForbidden,
			Detail:    error.Error(),
			Title:     "forbidden",
			Timestamp: time.Now(),
		}
	})

	problem.Map(http.StatusForbidden, func() *problem.ProblemDetail {
		return &problem.ProblemDetail{
			Type:      "https://httpstatuses.io/401",
			Status:    http.StatusUnauthorized,
			Detail:    error.Error(),
			Title:     "unauthorized",
			Timestamp: time.Now(),
		}
	})

	problem.Map(http.StatusForbidden, func() *problem.ProblemDetail {
		return &problem.ProblemDetail{
			Type:      "https://httpstatuses.io/404",
			Status:    http.StatusNotFound,
			Detail:    error.Error(),
			Title:     "not-found",
			Timestamp: time.Now(),
		}
	})

	if !c.Response().Committed {
		if err := problem.ResolveProblemDetails(c.Response(), error); err != nil {
			log.Error(err)
		}
	}
}
