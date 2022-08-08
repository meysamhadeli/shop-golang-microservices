package delivery

import (
	"github.com/labstack/echo/v4"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/shared/configurations/infrastructure"
)

type ProductEndpointBase struct {
	*infrastructure.InfrastructureConfiguration
	ProductsGroup *echo.Group
}
