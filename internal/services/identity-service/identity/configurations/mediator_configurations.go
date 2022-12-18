package configurations

import (
	"github.com/mehdihadeli/go-mediatr"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/identity-service/identity/dtos"
	v1 "github.com/meysamhadeli/shop-golang-microservices/internal/services/identity-service/identity/features/registering_user/commands/v1"
	contracts2 "github.com/meysamhadeli/shop-golang-microservices/internal/services/identity-service/shared/contracts"
)

func ConfigUsersMediator(ic *contracts2.InfrastructureConfiguration) error {

	//https://stackoverflow.com/questions/72034479/how-to-implement-generic-interfaces
	err := mediatr.RegisterRequestHandler[*v1.RegisterUser, *dtos.RegisterUserResponseDto](v1.NewRegisterUserHandler(ic.Log, ic.Cfg, ic.UserRepository, ic.RabbitmqPublisher))
	if err != nil {
		return err
	}

	return nil
}
