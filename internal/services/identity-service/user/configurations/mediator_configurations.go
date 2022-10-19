package configurations

import (
	"github.com/mehdihadeli/go-mediatr"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/identity-service/user/contracts"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/identity-service/user/dtos"
	v1 "github.com/meysamhadeli/shop-golang-microservices/internal/services/identity-service/user/features/registering_user/commands/v1"
)

func (c *usersModuleConfigurator) configUsersMediator(pgRepo contracts.UserRepository) error {

	//https://stackoverflow.com/questions/72034479/how-to-implement-generic-interfaces
	err := mediatr.RegisterRequestHandler[*v1.RegisterUser, *dtos.RegisterUserResponseDto](v1.NewRegisterUserHandler(c.Log, c.Cfg, pgRepo, c.RabbitmqPublisher, c.JaegerTracer))
	if err != nil {
		return err
	}

	return nil
}
