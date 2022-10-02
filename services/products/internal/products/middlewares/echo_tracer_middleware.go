package middlewares

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// ref: https://github.com/open-telemetry/opentelemetry-go-contrib/blob/main/instrumentation/github.com/labstack/echo/otelecho/echo.go

func EchoTracerMiddleware(serviceName string) echo.MiddlewareFunc {

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			request := c.Request()
			ctx := request.Context()
			defer func() {
				request = request.WithContext(ctx)
				c.SetRequest(request)
			}()

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

			// serve the request to the next middleware
			err := next(c)

			// invokes the registered HTTP error handler
			c.Error(err)

			var echoError *echo.HTTPError

			// handle *HTTPError error type in Echo
			if errors.As(err, &echoError) {
				c.Response().Status = err.(*echo.HTTPError).Code
				err = err.(*echo.HTTPError).Message.(error)
			}

			if err != nil {
				span.SetAttributes(attribute.String("echo.error", err.Error()))
			}

			attrs := semconv.HTTPAttributesFromHTTPStatusCode(c.Response().Status)
			spanStatus, spanMessage := semconv.SpanStatusFromHTTPStatusCode(c.Response().Status) // set the spanStatus Error for all error stats codes
			span.SetAttributes(attrs...)
			span.SetStatus(spanStatus, spanMessage)

			return nil
		}
	}
}
