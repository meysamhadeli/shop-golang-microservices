package config

import (
	"fmt"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/config"
	"net/url"
)

func Address(cfg config.Config) string {
	return fmt.Sprintf("%s%s", cfg.Echo.Host, cfg.Echo.Port)
}

func BasePathAddress(cfg config.Config) string {
	path, err := url.JoinPath(Address(cfg), cfg.Echo.BasePath)
	if err != nil {
		return ""
	}
	return path
}
