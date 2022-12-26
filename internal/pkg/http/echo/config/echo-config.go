package config

import (
	"fmt"
	"net/url"
)

type EchoConfig struct {
	Port                string   `mapstructure:"port" validate:"required"`
	Development         bool     `mapstructure:"development"`
	BasePath            string   `mapstructure:"basePath" validate:"required"`
	DebugErrorsResponse bool     `mapstructure:"debugErrorsResponse"`
	IgnoreLogUrls       []string `mapstructure:"ignoreLogUrls"`
	Timeout             int      `mapstructure:"timeout"`
	Host                string   `mapstructure:"host"`
}

func Address(cfg *EchoConfig) string {
	return fmt.Sprintf("%s%s", cfg.Host, cfg.Port)
}

func BasePathAddress(cfg *EchoConfig) string {
	path, err := url.JoinPath(Address(cfg), cfg.BasePath)
	if err != nil {
		return ""
	}
	return path
}
