package config

import (
	"flag"
	"fmt"
	"github.com/meysamhadeli/shop-golang-microservices/pkg/gorm_postgres"
	grpc "github.com/meysamhadeli/shop-golang-microservices/pkg/grpc/config"
	echo "github.com/meysamhadeli/shop-golang-microservices/pkg/http/echo/config"
	"github.com/meysamhadeli/shop-golang-microservices/pkg/logger"
	open_telemetry "github.com/meysamhadeli/shop-golang-microservices/pkg/open-telemetry"
	"github.com/meysamhadeli/shop-golang-microservices/pkg/rabbitmq"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/constants"
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
	Echo         *echo.EchoConfig         `mapstructure:"echo"`
	Grpc         *grpc.GrpcConfig         `mapstructure:"grpc"`
	Context      Context                  `mapstructure:"context"`
	GormPostgres *gorm_postgres.Config    `mapstructure:"gormPostgres"`
	Jaeger       *open_telemetry.Config   `mapstructure:"jaeger"`
}

type Context struct {
	Timeout int `mapstructure:"timeout"`
}

func InitConfig(env string) (*Config, error) {
	if configPath == "" {
		configPathFromEnv := os.Getenv(constants.ConfigPath)
		if configPathFromEnv != "" {
			configPath = configPathFromEnv
		} else {
			configPath = "./services/products/config"
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

func GetMicroserviceName(serviceName string) string {
	return fmt.Sprintf("%s", strings.ToUpper(serviceName))
}
