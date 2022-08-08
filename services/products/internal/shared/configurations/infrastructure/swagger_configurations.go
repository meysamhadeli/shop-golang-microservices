package infrastructure

import (
	echoSwagger "github.com/swaggo/echo-swagger"
)

func (ic *infrastructureConfigurator) configSwagger() {
	ic.echo.GET("/swagger/*", echoSwagger.WrapHandler)
}
