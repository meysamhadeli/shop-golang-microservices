package open_telemetry

import (
	"context"
	"encoding/json"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func TraceErr(ctx context.Context, tracer trace.Tracer, err error) {
	_, span := tracer.Start(ctx, "tracer_error")
	defer span.End()

	span.SetStatus(codes.Error, err.Error())
	attribute.Bool("error", true)
}

func TraceWithErr(ctx context.Context, tracer trace.Tracer, err error) error {
	_, span := tracer.Start(ctx, "tracer_error")
	defer span.End()

	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		attribute.Bool("error", true)
	}
	return err
}

func ObjToString(obj ...interface{}) (string, error) {
	value, err := json.Marshal(obj)
	if err != nil {
		return *new(string), err
	}
	return string(value), nil
}
