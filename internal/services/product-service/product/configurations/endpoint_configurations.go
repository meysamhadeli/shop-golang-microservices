package configurations

import (
	creating_product "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/creating_product/endpoints/v1"
	deleting_product "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/deleting_product/endpoints/v1"
	getting_product_by_id "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/getting_product_by_id/endpoints/v1"
	getting_products "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/getting_products/endpoints/v1"
	searching_product "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/searching_product/endpoints/v1"
	updating_product "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/updating_product/endpoints/v1"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/shared/contracts"
)

func ConfigEndpoints(infra *contracts.InfrastructureConfiguration) {

	creating_product.MapRoute(infra)
	deleting_product.MapRoute(infra)
	getting_product_by_id.MapRoute(infra)
	getting_products.MapRoute(infra)
	searching_product.MapRoute(infra)
	updating_product.MapRoute(infra)
}
