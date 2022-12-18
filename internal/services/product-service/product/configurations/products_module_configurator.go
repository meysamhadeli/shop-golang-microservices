package configurations

import (
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/grpc"
	repositories_imp "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/data/repositories"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/mappings"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/shared/contracts"
)

type ProductsModuleConfigurator interface {
	ConfigureProductsModule() error
}

func NewProductsModuleConfigurator(infrastructure *contracts.InfrastructureConfiguration, identityGrpcClient grpc.GrpcClient) *contracts.InfrastructureConfiguration {
	return infrastructure
}

func ConfigureProductsModule(ic *contracts.InfrastructureConfiguration) error {

	v1 := ic.Echo.Group("/api/v1")
	group := v1.Group("/products")

	ic.ProductRepository = repositories_imp.NewPostgresProductRepository(ic.Log, ic.Cfg, ic.Gorm)

	err := mappings.ConfigureMappings()
	if err != nil {
		return err
	}

	err = ConfigProductsMediator(ic)
	if err != nil {
		return err
	}

	ConfigEndpoints(ic, group)

	return nil
}
