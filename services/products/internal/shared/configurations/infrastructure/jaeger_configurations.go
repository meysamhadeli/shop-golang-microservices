package infrastructure

import (
	"github.com/meysamhadeli/shop-golang-microservices/pkg/tracing"
	"github.com/opentracing/opentracing-go"
)

func (ic *infrastructureConfigurator) configJaeger() (error, func()) {
	if ic.cfg.Jaeger.Enable {
		tracer, closer, err := tracing.NewJaegerTracer(ic.cfg.Jaeger)
		if err != nil {
			return err, nil
		}
		opentracing.SetGlobalTracer(tracer)
		return nil, func() {
			_ = closer.Close()
		}
	}

	return nil, func() {}
}
