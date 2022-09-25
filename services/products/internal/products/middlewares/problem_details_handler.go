package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/meysamhadeli/shop-golang-microservices/pkg/problemDetails"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func ProblemDetailsHandler(error error, c echo.Context) {

	problemDetails.Map(http.StatusBadRequest, func() *problemDetails.ProblemDetail {
		return &problemDetails.ProblemDetail{
			Type:      "https://httpstatuses.io/400",
			Detail:    error.Error(),
			Title:     "bad-request",
			Timestamp: time.Now(),
		}
	})

	problemDetails.Map(http.StatusForbidden, func() *problemDetails.ProblemDetail {
		return &problemDetails.ProblemDetail{
			Type:      "https://httpstatuses.io/403",
			Status:    http.StatusForbidden,
			Detail:    error.Error(),
			Title:     "forbidden",
			Timestamp: time.Now(),
		}
	})

	problemDetails.Map(http.StatusForbidden, func() *problemDetails.ProblemDetail {
		return &problemDetails.ProblemDetail{
			Type:      "https://httpstatuses.io/401",
			Status:    http.StatusUnauthorized,
			Detail:    error.Error(),
			Title:     "unauthorized",
			Timestamp: time.Now(),
		}
	})

	problemDetails.Map(http.StatusForbidden, func() *problemDetails.ProblemDetail {
		return &problemDetails.ProblemDetail{
			Type:      "https://httpstatuses.io/404",
			Status:    http.StatusNotFound,
			Detail:    error.Error(),
			Title:     "not-fund",
			Timestamp: time.Now(),
		}
	})

	if !c.Response().Committed {
		if _, err := problemDetails.ResolveProblemDetails(c.Response(), error); err != nil {
			log.Error(err)
		}
	}
}
