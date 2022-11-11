package configurations

import (
	"context"
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

func NewUsersModuleConfigurator(infrastructure *contracts.InfrastructureConfiguration) *usersModuleConfigurator {
	return &usersModuleConfigurator{InfrastructureConfiguration: infrastructure}
}

func (c *usersModuleConfigurator) ConfigureIdentitiesModule(ctx context.Context) error {

	v1 := c.Echo.Group("/api/v1")
	group := v1.Group("/users")

	userRepository := repositories_imp.NewPostgresUserRepository(c.Log, c.Cfg, c.Gorm)

	err := mappings.ConfigureMappings()
	if err != nil {
		return err
	}

	err = c.configUsersMediator(userRepository)
	if err != nil {
		return err
	}

	c.configEndpoints(ctx, group)

	return nil
}
