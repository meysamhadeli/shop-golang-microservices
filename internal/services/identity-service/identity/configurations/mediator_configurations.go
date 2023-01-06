package configurations

import (
	"context"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/logger"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/rabbitmq"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/identity-service/identity/contracts"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/identity-service/identity/dtos"
	v1 "github.com/meysamhadeli/shop-golang-microservices/internal/services/identity-service/identity/features/registering_user/commands/v1"
)

func ConfigUsersMediator(log logger.ILogger, rabbitmqPublisher rabbitmq.IPublisher,
	userRepository contracts.UserRepository, ctx context.Context) error {

	//https://stackoverflow.com/questions/72034479/how-to-implement-generic-interfaces
	err := mediatr.RegisterRequestHandler[*v1.RegisterUser, *dtos.RegisterUserResponseDto](v1.NewRegisterUserHandler(log, rabbitmqPublisher, userRepository, ctx))
	if err != nil {
		return err
	}

	return nil
}
