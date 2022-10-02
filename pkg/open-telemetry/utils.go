package open_telemetry

import (
	"context"

	"go.opentelemetry.io/otel"
)

//ref: https://devandchill.com/posts/2021/12/go-step-by-step-guide-for-implementing-tracing-on-a-microservices-architecture-2/2/

type AmqpHeadersCarrier map[string]interface{}

func (a AmqpHeadersCarrier) Get(key string) string {
	v, ok := a[key]
	if !ok {
		return ""
	}
	return v.(string)
}

func (a AmqpHeadersCarrier) Set(key string, value string) {
	a[key] = value
}

func (a AmqpHeadersCarrier) Keys() []string {
	i := 0
	r := make([]string, len(a))

	for k, _ := range a {
		r[i] = k
		i++
	}

	return r
}

// InjectAMQPHeaders injects the tracing from the context into the header map
func InjectAMQPHeaders(ctx context.Context) map[string]interface{} {
	h := make(AmqpHeadersCarrier)
	otel.GetTextMapPropagator().Inject(ctx, h)
	return h
}

// ExtractAMQPHeaders extracts the tracing from the header and puts it into the context
func ExtractAMQPHeaders(ctx context.Context, headers map[string]interface{}) context.Context {
	return otel.GetTextMapPropagator().Extract(ctx, AmqpHeadersCarrier(headers))
}
