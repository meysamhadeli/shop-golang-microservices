package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/meysamhadeli/problem-details"
	log "github.com/sirupsen/logrus"
)

func ProblemDetailsHandler(error error, c echo.Context) {
	if !c.Response().Committed {
		if err := problem.ResolveProblemDetails(c.Response(), c.Request(), error); err != nil {
			log.Error(err)
		}
	}
}
