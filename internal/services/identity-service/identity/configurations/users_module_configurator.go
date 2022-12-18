package configurations

import (
	repositories_imp "github.com/meysamhadeli/shop-golang-microservices/internal/services/identity-service/identity/data/repositories"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/identity-service/identity/mappings"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/identity-service/shared/contracts"
)

type UsersModuleConfigurator interface {
	ConfigureUsersModule() error
}

type usersModuleConfigurator struct {
	*contracts.InfrastructureConfiguration
}

func NewUsersModuleConfigurator(infrastructure *contracts.InfrastructureConfiguration) *contracts.InfrastructureConfiguration {
	return infrastructure
}

func ConfigureIdentitiesModule(ic *contracts.InfrastructureConfiguration) error {

	v1 := ic.Echo.Group("/api/v1")
	group := v1.Group("/users")

	ic.UserRepository = repositories_imp.NewPostgresUserRepository(ic.Log, ic.Cfg, ic.Gorm)

	err := mappings.ConfigureMappings()
	if err != nil {
		return err
	}

	err = ConfigUsersMediator(ic)
	if err != nil {
		return err
	}

	ConfigEndpoints(ic, group)

	return nil
}
