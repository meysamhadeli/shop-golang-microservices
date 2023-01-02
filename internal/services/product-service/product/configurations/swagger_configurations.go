package configurations

import (
	echo_server "github.com/meysamhadeli/shop-golang-microservices/internal/pkg/http/echo/server"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/docs"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func ConfigSwagger(e *echo_server.EchoServer) {

	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Title = "Products Service Api"
	docs.SwaggerInfo.Description = "Products Service Api"
	e.Echo.GET("/swagger/*", echoSwagger.WrapHandler)
}
