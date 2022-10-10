package middleware

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

func EchoTracerMiddleware(serviceName string) echo.MiddlewareFunc {

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			request := c.Request()
			ctx := request.Context()

			// ref: https://github.com/open-telemetry/opentelemetry-go-contrib/blob/main/instrumentation/github.com/labstack/echo/otelecho/echo.go
			ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.HeaderCarrier(request.Header))
			opts := []oteltrace.SpanStartOption{
				oteltrace.WithAttributes(semconv.NetAttributesFromHTTPRequest("tcp", request)...),
				oteltrace.WithAttributes(semconv.EndUserAttributesFromHTTPRequest(request)...),
				oteltrace.WithAttributes(semconv.HTTPServerAttributesFromHTTPRequest(serviceName, c.Path(), request)...),
				oteltrace.WithSpanKind(oteltrace.SpanKindServer),
			}
			spanName := c.Path()
			if spanName == "" {
				spanName = fmt.Sprintf("HTTP %s route not found", request.Method)
			}

			ctx, span := otel.Tracer("echo-http").Start(ctx, spanName, opts...)
			defer span.End()

			// pass the span through the request context
			c.SetRequest(request.WithContext(ctx))

			err := next(c)

			if err != nil {
				// invokes the registered HTTP error handler
				c.Error(err)

				var echoError *echo.HTTPError

				// handle *HTTPError error type in Echo
				if errors.As(err, &echoError) {
					c.Response().Status = err.(*echo.HTTPError).Code
					err = err.(*echo.HTTPError).Message.(error)
				}

				span.SetStatus(codes.Error, "") // set the spanStatus Error for all error stats codes
				span.SetAttributes(attribute.String("echo-error", err.Error()))
				span.SetAttributes(attribute.Int("status-code", c.Response().Status))
			}

			return err
		}
	}
}
