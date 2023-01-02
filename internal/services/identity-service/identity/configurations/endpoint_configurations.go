package configurations

import (
	registering_user "github.com/meysamhadeli/shop-golang-microservices/internal/services/identity-service/identity/features/registering_user/endpoints/v1"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/identity-service/shared/contracts"
)

func ConfigEndpoints(infra *contracts.InfrastructureConfiguration) {

	registering_user.MapRoute(infra)
}
