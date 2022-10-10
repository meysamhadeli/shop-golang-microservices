package configurations

import (
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/docs"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func (ic *infrastructureConfigurator) configSwagger() {

	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Title = "Products Service Api"
	docs.SwaggerInfo.Description = "Products Service Api"
	ic.Echo.GET("/swagger/*", echoSwagger.WrapHandler)
}
