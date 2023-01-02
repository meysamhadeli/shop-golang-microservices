package configurations

import (
	identity_service "github.com/meysamhadeli/shop-golang-microservices/internal/services/identity-service/identity/grpc_server/protos"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/identity-service/identity/grpc_server/services"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/identity-service/shared/contracts"
)

func ConfigIdentityGrpcServer(infra *contracts.InfrastructureConfiguration) {

	identityGrpcService := services.NewIdentityGrpcServerService(infra)

	identity_service.RegisterIdentityServiceServer(infra.GrpcServer.Grpc, identityGrpcService)
}
