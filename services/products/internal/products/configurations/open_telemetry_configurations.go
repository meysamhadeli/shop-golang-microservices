package configurations

import (
	open_telemetry "github.com/meysamhadeli/shop-golang-microservices/pkg/open-telemetry"
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
