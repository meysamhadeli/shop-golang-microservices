package v1

import (
	"context"
	"encoding/json"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/mapper"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/creating_product/dtos/v1"
	v12 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/creating_product/events/v1"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/models"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/shared/contracts"
)

type CreateProductHandler struct {
	infra *contracts.InfrastructureConfiguration
}

func NewCreateProductHandler(infra *contracts.InfrastructureConfiguration) *CreateProductHandler {
	return &CreateProductHandler{infra: infra}
}

func (c *CreateProductHandler) Handle(ctx context.Context, command *CreateProduct) (*v1.CreateProductResponseDto, error) {

	product := &models.Product{
		ProductId:   command.ProductID,
		Name:        command.Name,
		Description: command.Description,
		Price:       command.Price,
		CreatedAt:   command.CreatedAt,
	}

	createdProduct, err := c.infra.ProductRepository.CreateProduct(ctx, product)
	if err != nil {
		return nil, err
	}

	evt, err := mapper.Map[*v12.ProductCreated](createdProduct)
	if err != nil {
		return nil, err
	}

	err = c.infra.RabbitmqPublisher.PublishMessage(ctx, evt)
	if err != nil {
		return nil, err
	}

	response := &v1.CreateProductResponseDto{ProductId: product.ProductId}
	bytes, _ := json.Marshal(response)

	c.infra.Log.Info("CreateProductResponseDto", string(bytes))

	return response, nil
}
