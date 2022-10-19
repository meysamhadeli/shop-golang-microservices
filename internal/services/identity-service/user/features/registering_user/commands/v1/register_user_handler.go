package v1

import (
	"context"
	"encoding/json"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/logger"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/mapper"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/rabbitmq"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/identity-service/config"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/identity-service/user/contracts"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/identity-service/user/dtos"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/identity-service/user/models"
	"go.opentelemetry.io/otel/trace"
)

type RegisterUserHandler struct {
	log               logger.ILogger
	cfg               *config.Config
	repository        contracts.UserRepository
	rabbitmqPublisher rabbitmq.IPublisher
	jaegerTracer      trace.Tracer
}

func NewRegisterUserHandler(log logger.ILogger, cfg *config.Config, repository contracts.UserRepository,
	rabbitmqPublisher rabbitmq.IPublisher, jaegerTracer trace.Tracer) *RegisterUserHandler {
	return &RegisterUserHandler{log: log, cfg: cfg, repository: repository, rabbitmqPublisher: rabbitmqPublisher, jaegerTracer: jaegerTracer}
}

func (c *RegisterUserHandler) Handle(ctx context.Context, command *RegisterUser) (*dtos.RegisterUserResponseDto, error) {

	product := &models.User{
		UserId:    command.UserId,
		Email:     command.Email,
		Password:  command.Password,
		UserName:  command.UserName,
		LastName:  command.LastName,
		FirstName: command.FirstName,
		CreatedAt: command.CreatedAt,
	}

	registeredUser, err := c.repository.RegisterUser(ctx, product)
	if err != nil {
		return nil, err
	}

	response, err := mapper.Map[*dtos.RegisterUserResponseDto](registeredUser)
	if err != nil {
		return nil, err
	}
	bytes, _ := json.Marshal(response)

	c.log.Info("RegisterUserResponseDto", string(bytes))

	return response, nil
}
