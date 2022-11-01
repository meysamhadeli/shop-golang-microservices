package services

import (
	"context"
	identity_service "github.com/meysamhadeli/shop-golang-microservices/internal/services/identity-service/identity/grpc_server/protos"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/identity-service/shared"
	uuid "github.com/satori/go.uuid"
)

type IdentityGrpcServerService struct {
	*shared.InfrastructureConfiguration
}

func NewIdentityGrpcServerService(infra *shared.InfrastructureConfiguration) *IdentityGrpcServerService {
	return &IdentityGrpcServerService{InfrastructureConfiguration: infra}
}

func (i IdentityGrpcServerService) GetUserById(ctx context.Context, req *identity_service.GetUserByIdReq) (*identity_service.GetUserByIdRes, error) {

	var user = &identity_service.User{UserId: uuid.NewV4().String(), Name: "sam"}

	var result = &identity_service.GetUserByIdRes{
		User: user,
	}

	return result, nil
}
