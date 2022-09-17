package config

import (
	"flag"
	"fmt"
	"github.com/meysamhadeli/shop-golang-microservices/pkg/constants"
	"github.com/meysamhadeli/shop-golang-microservices/pkg/gorm_postgres"
	"github.com/meysamhadeli/shop-golang-microservices/pkg/logger"
	open_telemetry "github.com/meysamhadeli/shop-golang-microservices/pkg/open-telemetry"
	"github.com/meysamhadeli/shop-golang-microservices/pkg/rabbitmq"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "", "products write microservice config path")
}

type Config struct {
	ServiceName  string                   `mapstructure:"serviceName"`
	Logger       *logger.Config           `mapstructure:"logger"`
	Rabbitmq     *rabbitmq.RabbitMQConfig `mapstructure:"rabbitmq"`
	Http         Http                     `mapstructure:"http"`
	Context      Context                  `mapstructure:"context"`
	GormPostgres *gorm_postgres.Config    `mapstructure:"gormPostgres"`
	Jaeger       *open_telemetry.Config   `mapstructure:"jaeger"`
}

type Context struct {
	Timeout int `mapstructure:"timeout"`
}

type Http struct {
	Port                string   `mapstructure:"port" validate:"required"`
	Development         bool     `mapstructure:"development"`
	BasePath            string   `mapstructure:"basePath" validate:"required"`
	ProductsPath        string   `mapstructure:"productsPath" validate:"required"`
	DebugErrorsResponse bool     `mapstructure:"debugErrorsResponse"`
	IgnoreLogUrls       []string `mapstructure:"ignoreLogUrls"`
	Timeout             int      `mapstructure:"timeout"`
	Host                string   `mapstructure:"host"`
}

func InitConfig(env string) (*Config, error) {
	if configPath == "" {
		configPathFromEnv := os.Getenv(constants.ConfigPath)
		if configPathFromEnv != "" {
			configPath = configPathFromEnv
		} else {
			configPath = "./config"
		}
	}

	cfg := &Config{}

	viper.SetConfigName(fmt.Sprintf("config.%s", env))
	viper.AddConfigPath(configPath)
	viper.SetConfigType(constants.Json)

	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "viper.ReadInConfig")
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, errors.Wrap(err, "viper.Unmarshal")
	}

	return cfg, nil
}

func GetMicroserviceName(cfg *Config) string {
	return fmt.Sprintf("(%s)", strings.ToUpper(cfg.ServiceName))
}
