package config

type GrpcConfig struct {
	Port        string `mapstructure:"port"`
	Host        string `mapstructure:"host"`
	Development bool   `mapstructure:"development"`
}
