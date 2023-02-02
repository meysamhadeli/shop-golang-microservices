package configurations

import (
	"context"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/logger"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/rabbitmq"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/identity-service/identity/contracts"
	registeringuserv1commands "github.com/meysamhadeli/shop-golang-microservices/internal/services/identity-service/identity/features/registering_user/v1/commands"
	registeringuserv1dtos "github.com/meysamhadeli/shop-golang-microservices/internal/services/identity-service/identity/features/registering_user/v1/dtos"
)

func ConfigUsersMediator(log logger.ILogger, rabbitmqPublisher rabbitmq.IPublisher,
	userRepository contracts.UserRepository, ctx context.Context) error {

	//https://stackoverflow.com/questions/72034479/how-to-implement-generic-interfaces
	err := mediatr.RegisterRequestHandler[*registeringuserv1commands.RegisterUser, *registeringuserv1dtos.RegisterUserResponseDto](registeringuserv1commands.NewRegisterUserHandler(log, rabbitmqPublisher, userRepository, ctx))
	if err != nil {
		return err
	}

	return nil
}
