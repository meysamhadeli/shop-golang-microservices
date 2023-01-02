package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/mapper"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/updating_product/dtos/v1"
	events_v1 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/updating_product/events/v1"
	identity_service "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/grpc_client/protos"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/models"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/shared/contracts"
	"github.com/pkg/errors"
)

type UpdateProductHandler struct {
	infra *contracts.InfrastructureConfiguration
}

func NewUpdateProductHandler(infra *contracts.InfrastructureConfiguration) *UpdateProductHandler {
	return &UpdateProductHandler{infra: infra}
}

func (c *UpdateProductHandler) Handle(ctx context.Context, command *UpdateProduct) (*v1.UpdateProductResponseDto, error) {

	//simple call grpcClient
	identityGrpcClient := identity_service.NewIdentityServiceClient(c.infra.GrpcClient.GetGrpcConnection())
	user, err := identityGrpcClient.GetUserById(ctx, &identity_service.GetUserByIdReq{UserId: "1"})
	if err != nil {
		return nil, err
	}

	c.infra.Log.Infof("userId: %s", user.User.UserId)

	_, err = c.infra.ProductRepository.GetProductById(ctx, command.ProductID)

	if err != nil {
		notFoundErr := errors.Wrap(err, fmt.Sprintf("product with id %s not found", command.ProductID))
		return nil, notFoundErr
	}

	product := &models.Product{ProductId: command.ProductID, Name: command.Name, Description: command.Description, Price: command.Price, UpdatedAt: command.UpdatedAt}

	updatedProduct, err := c.infra.ProductRepository.UpdateProduct(ctx, product)
	if err != nil {
		return nil, err
	}

	evt, err := mapper.Map[*events_v1.ProductUpdated](updatedProduct)
	if err != nil {
		return nil, err
	}

	err = c.infra.RabbitmqPublisher.PublishMessage(ctx, evt)

	response := &v1.UpdateProductResponseDto{ProductId: product.ProductId}
	bytes, _ := json.Marshal(response)

	c.infra.Log.Info("UpdateProductResponseDto", string(bytes))

	return response, nil
}
