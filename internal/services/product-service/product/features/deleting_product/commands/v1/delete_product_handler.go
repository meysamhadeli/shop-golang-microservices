package v1

import (
	"context"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/deleting_product/events"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/shared/contracts"
)

type DeleteProductHandler struct {
	infra *contracts.InfrastructureConfiguration
}

func NewDeleteProductHandler(infra *contracts.InfrastructureConfiguration) *DeleteProductHandler {
	return &DeleteProductHandler{infra: infra}
}

func (c *DeleteProductHandler) Handle(ctx context.Context, command *DeleteProduct) (*mediatr.Unit, error) {

	if err := c.infra.ProductRepository.DeleteProductByID(ctx, command.ProductID); err != nil {
		return nil, err
	}

	err := c.infra.RabbitmqPublisher.PublishMessage(ctx, events.ProductDeleted{
		ProductId: command.ProductID,
	})
	if err != nil {
		return nil, err
	}

	c.infra.Log.Info("DeleteProduct successfully executed")

	return &mediatr.Unit{}, err
}
