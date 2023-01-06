package configurations

import (
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/grpc"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/logger"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/identity-service/config"
	identity_service "github.com/meysamhadeli/shop-golang-microservices/internal/services/identity-service/identity/grpc_server/protos"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/identity-service/identity/grpc_server/services"
)

func ConfigIdentityGrpcServer(cfg *config.Config, log logger.ILogger, grpcServer *grpc.GrpcServer) {

	identityGrpcService := services.NewIdentityGrpcServerService(cfg, log)

	identity_service.RegisterIdentityServiceServer(grpcServer.Grpc, identityGrpcService)
}
