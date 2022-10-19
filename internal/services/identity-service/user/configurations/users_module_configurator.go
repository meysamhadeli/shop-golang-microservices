package configurations

import (
	"context"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/identity-service/shared"
	repositories_imp "github.com/meysamhadeli/shop-golang-microservices/internal/services/identity-service/user/data/repositories"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/identity-service/user/mappings"
)

type UsersModuleConfigurator interface {
	ConfigureUsersModule() error
}

type usersModuleConfigurator struct {
	*shared.InfrastructureConfiguration
}

func NewUsersModuleConfigurator(infrastructure *shared.InfrastructureConfiguration) *usersModuleConfigurator {
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
