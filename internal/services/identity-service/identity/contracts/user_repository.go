package contracts

import (
	"context"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/identity-service/identity/models"
)

type UserRepository interface {
	RegisterUser(ctx context.Context, user *models.User) (*models.User, error)
}
