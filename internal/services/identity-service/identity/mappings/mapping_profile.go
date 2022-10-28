package mappings

import (
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/mapper"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/identity-service/identity/dtos"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/identity-service/identity/models"
)

func ConfigureMappings() error {
	err := mapper.CreateMap[*models.User, *dtos.RegisterUserResponseDto]()
	if err != nil {
		return err
	}
	return err
}
