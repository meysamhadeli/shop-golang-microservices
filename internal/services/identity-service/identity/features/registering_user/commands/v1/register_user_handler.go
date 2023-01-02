package v1

import (
	"context"
	"encoding/json"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/mapper"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/utils"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/identity-service/identity/dtos"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/identity-service/identity/models"
	contracts2 "github.com/meysamhadeli/shop-golang-microservices/internal/services/identity-service/shared/contracts"
)

type RegisterUserHandler struct {
	infra *contracts2.InfrastructureConfiguration
}

func NewRegisterUserHandler(infra *contracts2.InfrastructureConfiguration) *RegisterUserHandler {
	return &RegisterUserHandler{infra: infra}
}

func (c *RegisterUserHandler) Handle(ctx context.Context, command *RegisterUser) (*dtos.RegisterUserResponseDto, error) {

	pass, err := utils.HashPassword(command.Password)
	if err != nil {
		return nil, err
	}

	product := &models.User{
		UserId:    command.UserId,
		Email:     command.Email,
		Password:  pass,
		UserName:  command.UserName,
		LastName:  command.LastName,
		FirstName: command.FirstName,
		CreatedAt: command.CreatedAt,
	}

	registeredUser, err := c.infra.UserRepository.RegisterUser(ctx, product)
	if err != nil {
		return nil, err
	}

	response, err := mapper.Map[*dtos.RegisterUserResponseDto](registeredUser)
	if err != nil {
		return nil, err
	}
	bytes, _ := json.Marshal(response)

	c.infra.Log.Info("RegisterUserResponseDto", string(bytes))

	return response, nil
}
