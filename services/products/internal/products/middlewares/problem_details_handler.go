package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/meysamhadeli/shop-golang-microservices/pkg/problemDetails"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func ProblemDetailsHandler(err error, c echo.Context) {

	problemDetails.Map(err, func() *problemDetails.ProblemDetail {
		return &problemDetails.ProblemDetail{
			Type:      "https://httpstatuses.io/409",
			Status:    http.StatusConflict,
			Detail:    err.Error(),
			Title:     "application rule broken",
			Timestamp: time.Now(),
		}
	})

	problemDetails.Map(err, func() *problemDetails.ProblemDetail {
		return &problemDetails.ProblemDetail{
			Type:      "https://httpstatuses.io/400",
			Status:    http.StatusBadRequest,
			Detail:    err.Error(),
			Title:     "application exception",
			Timestamp: time.Now(),
		}
	})

	if !c.Response().Committed {
		if _, err := problemDetails.ResolveEcho(c.Response(), err); err != nil {
			log.Error(err)
		}
	}
}
