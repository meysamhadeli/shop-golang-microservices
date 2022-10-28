package open_telemetry

import (
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"os"
)

type Config struct {
	Server      string `mapstructure:"server"`
	ServiceName string `mapstructure:"serviceName"`
	TracerName  string `mapstructure:"tracerName"`
}

func TracerProvider(cfg *Config) (*tracesdk.TracerProvider, error) {
	var serverUrl = fmt.Sprintf(cfg.Server+"%s", "/api/traces")
	// Create the Jaeger exporter
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(serverUrl)))
	if err != nil {
		return nil, err
	}

	env := os.Getenv("APP_ENV")

	if env != "production" {
		env = "development"
	}

	tp := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exporter),
		// Record information about this application in a Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(cfg.ServiceName),
			attribute.String("environment", env),
		)),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}))

	return tp, nil
}
