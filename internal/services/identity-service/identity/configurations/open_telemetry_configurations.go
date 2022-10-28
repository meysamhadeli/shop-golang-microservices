package configurations

import (
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/open-telemetry"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

func (ic *infrastructureConfigurator) configOpenTelemetry() (*tracesdk.TracerProvider, error) {

	tp, err := open_telemetry.TracerProvider(ic.Cfg.Jaeger)
	if err != nil {
		ic.Log.Fatal(err)
		return nil, err
	}

	return tp, nil
}
