package context_provider

import (
	"context"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/logger"
	"os"
	"os/signal"
	"syscall"
)

func NewContext(log logger.ILogger) context.Context {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Info("context is canceled!")
				cancel()
				return
			}
		}
	}()

	return ctx
}
