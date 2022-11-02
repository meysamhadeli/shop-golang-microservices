package configurations

import (
	"context"
	identity_service "github.com/meysamhadeli/shop-golang-microservices/internal/services/identity-service/identity/grpc_server/protos"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/identity-service/identity/grpc_server/services"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/identity-service/shared"
	"google.golang.org/grpc"
)

func ConfigIdentityGrpcServer(ctx context.Context, server *grpc.Server, infra *shared.InfrastructureConfiguration) {

	identityGrpcService := services.NewIdentityGrpcServerService(infra)

	identity_service.RegisterIdentityServiceServer(server, identityGrpcService)
}
