package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/meysamhadeli/shop-golang-microservices/pkg/problemDetails"
	httpErrors "github.com/meysamhadeli/shop-golang-microservices/pkg/problemDetails/custome_error/utils"
	log "github.com/sirupsen/logrus"
)

func ProblemDetailsHandler(err error, c echo.Context) {
	prb := problemDetails.ParseError(err)

	if prb != nil {
		if !c.Response().Committed {
			if _, err := prb.WriteTo(c.Response()); err != nil {
				log.Error(err)
			}
		}
	} else {
		if !c.Response().Committed {
			prb := problemDetails.NewInternalServerProblemDetail(err.Error(), httpErrors.ErrorsWithStack(err))
			if _, err := prb.WriteTo(c.Response()); err != nil {
				log.Error(err)
			}
		}
	}
}
