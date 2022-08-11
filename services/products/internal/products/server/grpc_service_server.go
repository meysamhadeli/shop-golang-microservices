package server

import (
	"context"
	"github.com/meysamhadeli/shop-golang-microservices/pkg/mapper"
	"github.com/meysamhadeli/shop-golang-microservices/pkg/mediatr"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/config"
	getting_product_by_id_dtos "github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/features/getting_product_by_id/dtos"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/models"

	product_service_client "github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/contracts/grpc/service_clients"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/features/creating_product"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/features/creating_product/dtos"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/features/getting_product_by_id"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/features/updating_product"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/mappings"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProductGrpcServiceServer struct {
	infrastructure *config.InfrastructureConfiguration
	// Ref:https://github.com/grpc/grpc-go/issues/3794#issuecomment-720599532
	// product_service_client.UnimplementedProductsServiceServer
}

func NewProductGrpcService(infra *config.InfrastructureConfiguration) *ProductGrpcServiceServer {
	return &ProductGrpcServiceServer{infrastructure: infra}
}

func (s *ProductGrpcServiceServer) CreateProduct(ctx context.Context, req *product_service_client.CreateProductReq) (*product_service_client.CreateProductRes, error) {

	command := creating_product.NewCreateProduct(req.GetName(), req.GetDescription(), req.GetPrice())

	if err := s.infrastructure.Validator.StructCtx(ctx, command); err != nil {
		s.infrastructure.Log.Errorf("(validate) err: {%v}", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}

	result, err := mediatr.Send[*dtos.CreateProductResponseDto](ctx, command)
	if err != nil {
		s.infrastructure.Log.Errorf("(CreateProduct.Handle) productId: {%s}, err: {%v}", command.ProductID, err)
		return nil, s.errResponse(codes.Internal, err)
	}

	s.infrastructure.Log.Infof("(product created) productId: {%s}", command.ProductID)

	return &product_service_client.CreateProductRes{ProductID: result.ProductID.String()}, nil
}

func (s *ProductGrpcServiceServer) UpdateProduct(ctx context.Context, req *product_service_client.UpdateProductReq) (*product_service_client.UpdateProductRes, error) {

	productUUID, err := uuid.FromString(req.GetProductID())
	if err != nil {
		s.infrastructure.Log.Warn("uuid.FromString", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}

	command := updating_product.NewUpdateProduct(productUUID, req.GetName(), req.GetDescription(), req.GetPrice())

	if err := s.infrastructure.Validator.StructCtx(ctx, command); err != nil {
		s.infrastructure.Log.Warn("validate", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}

	_, err = mediatr.Send[*mediatr.Unit](ctx, command)
	if err != nil {
		s.infrastructure.Log.Warn("UpdateProduct.Handle", err)
		return nil, s.errResponse(codes.Internal, err)
	}

	s.infrastructure.Log.Infof("(product updated) id: {%s}", productUUID.String())

	return &product_service_client.UpdateProductRes{}, nil
}

func (s *ProductGrpcServiceServer) GetProductById(ctx context.Context, req *product_service_client.GetProductByIdReq) (*product_service_client.GetProductByIdRes, error) {

	productUUID, err := uuid.FromString(req.GetProductID())
	if err != nil {
		return nil, s.errResponse(codes.InvalidArgument, err)
	}

	query := getting_product_by_id.NewGetProductById(productUUID)
	if err := s.infrastructure.Validator.StructCtx(ctx, query); err != nil {
		s.infrastructure.Log.Warn("validate", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}

	queryResult, err := mediatr.Send[*getting_product_by_id_dtos.GetProductByIdResponseDto](ctx, query)
	if err != nil {
		s.infrastructure.Log.Warn("GetProductById.Handle", err)
		return nil, s.errResponse(codes.Internal, err)
	}

	product, err := mapper.Map[*models.Product](queryResult.Product)

	if err != nil {
		return nil, s.errResponse(codes.Internal, err)
	}

	return &product_service_client.GetProductByIdRes{Product: mappings.WriterProductToGrpc(product)}, nil
}

func (s *ProductGrpcServiceServer) errResponse(c codes.Code, err error) error {
	return status.Error(c, err.Error())
}
