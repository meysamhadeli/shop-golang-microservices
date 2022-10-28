package config

type GrpcConfig struct {
	Port        string `mapstructure:"port" validate:"required"`
	Development bool   `mapstructure:"development"`
}
