package commands

import (
	"context"
	"encoding/json"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/logger"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/mapper"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/rabbitmq"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/utils"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/identity_service/identity/data/contracts"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/identity_service/identity/features/registering_user/v1/dtos"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/identity_service/identity/models"
)

type RegisterUserHandler struct {
	log               logger.ILogger
	rabbitmqPublisher rabbitmq.IPublisher
	userRepository    contracts.UserRepository
	ctx               context.Context
}

func NewRegisterUserHandler(log logger.ILogger, rabbitmqPublisher rabbitmq.IPublisher,
	userRepository contracts.UserRepository, ctx context.Context) *RegisterUserHandler {
	return &RegisterUserHandler{log: log, userRepository: userRepository, ctx: ctx, rabbitmqPublisher: rabbitmqPublisher}
}

func (c *RegisterUserHandler) Handle(ctx context.Context, command *RegisterUser) (*dtos.RegisterUserResponseDto, error) {

	pass, err := utils.HashPassword(command.Password)
	if err != nil {
		return nil, err
	}

	product := &models.User{
		Email:     command.Email,
		Password:  pass,
		UserName:  command.UserName,
		LastName:  command.LastName,
		FirstName: command.FirstName,
		CreatedAt: command.CreatedAt,
	}

	registeredUser, err := c.userRepository.RegisterUser(ctx, product)
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
