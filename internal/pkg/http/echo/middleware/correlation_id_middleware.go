package echomiddleware

import (
	"context"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
)

func CorrelationIdMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		req := c.Request()

		id := req.Header.Get(echo.HeaderXCorrelationID)
		if id == "" {
			id = uuid.NewV4().String()
		}

		c.Response().Header().Set(echo.HeaderXCorrelationID, id)
		newReq := req.WithContext(context.WithValue(req.Context(), echo.HeaderXCorrelationID, id))
		c.SetRequest(newReq)

		return next(c)
	}
}
