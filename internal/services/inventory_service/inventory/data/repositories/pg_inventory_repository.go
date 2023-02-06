package repositories

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	gormpgsql "github.com/meysamhadeli/shop-golang-microservices/internal/pkg/gorm_pgsql"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/logger"
	contracts "github.com/meysamhadeli/shop-golang-microservices/internal/services/inventory_service/inventory/data/contracts"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/inventory_service/inventory/models"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type PostgresInventoryRepository struct {
	log  logger.ILogger
	cfg  *gormpgsql.GormPostgresConfig
	db   *pgxpool.Pool
	gorm *gorm.DB
}

func NewPostgresInventoryRepository(log logger.ILogger, cfg *gormpgsql.GormPostgresConfig, gorm *gorm.DB) contracts.InventoryRepository {
	return &PostgresInventoryRepository{log: log, cfg: cfg, gorm: gorm}
}

func (p *PostgresInventoryRepository) AddProductItemToInventory(ctx context.Context, productItem *models.ProductItem) (*models.ProductItem, error) {

	if err := p.gorm.Create(&productItem).Error; err != nil {
		return nil, errors.Wrap(err, "error in the inserting product into the database.")
	}

	return productItem, nil
}

func (p *PostgresInventoryRepository) GetProductInInventories(ctx context.Context, uuid uuid.UUID) (*models.ProductItem, error) {
	var productItem models.ProductItem

	if err := p.gorm.First(&productItem, uuid).Error; err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("can't find the product item with id %s into the database.", uuid))
	}

	return &productItem, nil
}
