package configurations

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/identity-service/shared"
	v1 "github.com/meysamhadeli/shop-golang-microservices/internal/services/identity-service/user/features/registering_user/endpoints/v1"
)

func (c *usersModuleConfigurator) configEndpoints(ctx context.Context, group *echo.Group) {

	userEndpointBase := &shared.IdentityEndpointBase[shared.InfrastructureConfiguration]{
		ProductsGroup: group,
		Configuration: *c.InfrastructureConfiguration,
	}

	registerUserEndpoint := v1.NewCreteUserEndpoint(userEndpointBase)
	registerUserEndpoint.MapRoute()
}
