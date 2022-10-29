package configurations

import (
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/oauth2"
)

func (ic *infrastructureConfigurator) configureOauth2() {
	oauth2.RunOauthServer(ic.Echo)
}
