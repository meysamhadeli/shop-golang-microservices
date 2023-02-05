package configurations

import (
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func ConfigSwagger(e *echo.Echo) {

	//docs.SwaggerInfo.Version = "1.0"
	//docs.SwaggerInfo.Title = "Products Service Api"
	//docs.SwaggerInfo.Description = "Products Service Api"
	e.GET("/swagger/*", echoSwagger.WrapHandler)
}
