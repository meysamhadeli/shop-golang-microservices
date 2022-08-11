package configurations

import (
	echoSwagger "github.com/swaggo/echo-swagger"
)

func (ic *infrastructureConfigurator) configSwagger() {
	ic.Echo.GET("/swagger/*", echoSwagger.WrapHandler)
}
