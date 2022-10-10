package config

type GrpcConfig struct {
	Port        string `mapstructure:"port" validate:"required"`
	Host        string `mapstructure:"host" validate:"required"`
	Development bool   `mapstructure:"development"`
}
