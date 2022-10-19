package v1

import (
	"context"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/logger"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/rabbitmq"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/config"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/contracts"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/deleting_product/events"
)

type DeleteProductHandler struct {
	log               logger.ILogger
	cfg               *config.Config
	pgRepo            contracts.ProductRepository
	rabbitmqPublisher rabbitmq.IPublisher
}

func NewDeleteProductHandler(log logger.ILogger, cfg *config.Config, pgRepo contracts.ProductRepository, rabbitmqPublisher rabbitmq.IPublisher) *DeleteProductHandler {
	return &DeleteProductHandler{log: log, cfg: cfg, pgRepo: pgRepo, rabbitmqPublisher: rabbitmqPublisher}
}

func (c *DeleteProductHandler) Handle(ctx context.Context, command *DeleteProduct) (*mediatr.Unit, error) {

	if err := c.pgRepo.DeleteProductByID(ctx, command.ProductID); err != nil {
		return nil, err
	}

	err := c.rabbitmqPublisher.PublishMessage(ctx, events.ProductDeleted{
		ProductId: command.ProductID,
	})
	if err != nil {
		return nil, err
	}

	c.log.Info("DeleteProduct successfully executed")

	return &mediatr.Unit{}, err
}
