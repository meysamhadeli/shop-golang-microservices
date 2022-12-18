package configurations

import (
	"github.com/labstack/echo/v4"
	v1 "github.com/meysamhadeli/shop-golang-microservices/internal/services/identity-service/identity/features/registering_user/endpoints/v1"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/identity-service/shared/contracts"
)

func ConfigEndpoints(ic *contracts.InfrastructureConfiguration, group *echo.Group) {

	userEndpointBase := &contracts.IdentityEndpointBase[contracts.InfrastructureConfiguration]{
		ProductsGroup: group,
		Configuration: *ic,
	}

	registerUserEndpoint := v1.NewCreteUserEndpoint(userEndpointBase)
	registerUserEndpoint.MapRoute()
}
