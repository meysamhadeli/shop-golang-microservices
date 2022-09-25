package updating_product

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/meysamhadeli/shop-golang-microservices/pkg/logger"
	"github.com/meysamhadeli/shop-golang-microservices/pkg/mapper"
	"github.com/meysamhadeli/shop-golang-microservices/pkg/rabbitmq"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/config"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/contracts"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/dtos"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/events"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/models"
	"github.com/pkg/errors"
)

type UpdateProductHandler struct {
	log               logger.ILogger
	cfg               *config.Config
	pgRepo            contracts.ProductRepository
	rabbitmqPublisher rabbitmq.IPublisher
}

func NewUpdateProductHandler(log logger.ILogger, cfg *config.Config, pgRepo contracts.ProductRepository,
	rabbitmqPublisher rabbitmq.IPublisher) *UpdateProductHandler {
	return &UpdateProductHandler{log: log, cfg: cfg, pgRepo: pgRepo, rabbitmqPublisher: rabbitmqPublisher}
}

func (c *UpdateProductHandler) Handle(ctx context.Context, command *UpdateProduct) (*dtos.UpdateProductResponseDto, error) {

	_, err := c.pgRepo.GetProductById(ctx, command.ProductID)

	if err != nil {
		notFoundErr := errors.Wrap(err, fmt.Sprintf("product with id %s not found", command.ProductID))
		return nil, notFoundErr
	}

	product := &models.Product{ProductId: command.ProductID, Name: command.Name, Description: command.Description, Price: command.Price, UpdatedAt: command.UpdatedAt}

	updatedProduct, err := c.pgRepo.UpdateProduct(ctx, product)
	if err != nil {
		return nil, err
	}

	evt, err := mapper.Map[*events.ProductUpdated](updatedProduct)
	if err != nil {
		return nil, err
	}

	err = c.rabbitmqPublisher.PublishMessage(ctx, evt)

	response := &dtos.UpdateProductResponseDto{ProductID: product.ProductId}
	bytes, _ := json.Marshal(response)

	c.log.Info("UpdateProductResponseDto", string(bytes))

	return response, nil
}
