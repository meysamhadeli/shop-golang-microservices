package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/logger"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/mapper"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/rabbitmq"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/config"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/contracts"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/updating_product/dtos/v1"
	events_v1 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/updating_product/events/v1"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/models"
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

func (c *UpdateProductHandler) Handle(ctx context.Context, command *UpdateProduct) (*v1.UpdateProductResponseDto, error) {

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

	evt, err := mapper.Map[*events_v1.ProductUpdated](updatedProduct)
	if err != nil {
		return nil, err
	}

	err = c.rabbitmqPublisher.PublishMessage(ctx, evt)

	response := &v1.UpdateProductResponseDto{ProductId: product.ProductId}
	bytes, _ := json.Marshal(response)

	c.log.Info("UpdateProductResponseDto", string(bytes))

	return response, nil
}
