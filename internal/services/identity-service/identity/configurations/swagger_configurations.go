package configurations

import (
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/identity-service/docs"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func (ic *infrastructureConfigurator) configSwagger() {
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Title = "Identities Service Api"
	docs.SwaggerInfo.Description = "Identities Service Api"
	ic.Echo.GET("/swagger/*", echoSwagger.WrapHandler)
}
