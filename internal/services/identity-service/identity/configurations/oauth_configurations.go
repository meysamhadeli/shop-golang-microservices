package configurations

import (
	"github.com/labstack/echo/v4"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/oauth2"
)

func configureOauth2(e *echo.Echo) {
	oauth2.RunOauthServer(e)
}
