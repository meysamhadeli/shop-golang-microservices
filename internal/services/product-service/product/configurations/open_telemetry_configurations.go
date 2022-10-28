package configurations

import (
	"context"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/open-telemetry"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

func (ic *infrastructureConfigurator) configOpenTelemetry(ctx context.Context) (*tracesdk.TracerProvider, error) {

	tp, err := open_telemetry.TracerProvider(ic.Cfg.Jaeger)
	if err != nil {
		ic.Log.Fatal(err)
		return nil, err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				err = tp.Shutdown(ctx)
				ic.Log.Error("open-telemetry exited properly")
				if err != nil {
					return
				}
			}
		}
	}()

	return tp, nil
}
