package configurations

import (
	"context"

	product_service "github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/contracts/grpc/service_clients"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/delivery/grpc"
)

func (c *productsModuleConfigurator) configGrpc(ctx context.Context) {
	productGrpcService := grpc.NewProductGrpcService(c.InfrastructureConfiguration)
	product_service.RegisterProductsServiceServer(c.GrpcServer, productGrpcService)
}
