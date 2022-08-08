package web

import (
	"fmt"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/config"
	"strings"
)

func GetMicroserviceName(cfg *config.Config) string {
	return fmt.Sprintf("(%s)", strings.ToUpper(cfg.ServiceName))
}
