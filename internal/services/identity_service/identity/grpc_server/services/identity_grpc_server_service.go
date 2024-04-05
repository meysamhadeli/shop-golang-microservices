package services

import (
	"context"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/logger"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/identity_service/config"
	identity_service "github.com/meysamhadeli/shop-golang-microservices/internal/services/identity_service/identity/grpc_server/protos"
)

type IdentityGrpcServerService struct {
	cfg *config.Config
	log logger.ILogger
}

func NewIdentityGrpcServerService(cfg *config.Config, log logger.ILogger) *IdentityGrpcServerService {
	return &IdentityGrpcServerService{log: log, cfg: cfg}
}

func (i IdentityGrpcServerService) GetUserById(ctx context.Context, req *identity_service.GetUserByIdReq) (*identity_service.GetUserByIdRes, error) {

	var user = &identity_service.User{UserId: req.UserId}

	var result = &identity_service.GetUserByIdRes{
		User: user,
	}

	return result, nil
}
