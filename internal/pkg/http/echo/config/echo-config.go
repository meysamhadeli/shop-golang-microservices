package config

import (
	"fmt"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/config_options"
	"net/url"
)

func Address(cfg config_options.Config) string {
	return fmt.Sprintf("%s%s", cfg.Echo.Host, cfg.Echo.Port)
}

func BasePathAddress(cfg config_options.Config) string {
	path, err := url.JoinPath(Address(cfg), cfg.Echo.BasePath)
	if err != nil {
		return ""
	}
	return path
}
