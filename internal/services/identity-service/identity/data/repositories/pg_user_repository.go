package repositories

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/logger"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/identity-service/config"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/identity-service/identity/contracts"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/identity-service/identity/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type PostgresUserRepository struct {
	log  logger.ILogger
	cfg  *config.Config
	db   *pgxpool.Pool
	gorm *gorm.DB
}

func NewPostgresUserRepository(log logger.ILogger, cfg *config.Config, gorm *gorm.DB) contracts.UserRepository {
	return PostgresUserRepository{log: log, cfg: cfg, gorm: gorm}
}

func (p PostgresUserRepository) RegisterUser(ctx context.Context, user *models.User) (*models.User, error) {

	if err := p.gorm.Create(&user).Error; err != nil {
		return nil, errors.Wrap(err, "error in the inserting user into the database.")
	}

	return user, nil
}
