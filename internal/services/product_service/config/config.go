package config

import (
	"flag"
	"fmt"
	gormpgsql "github.com/meysamhadeli/shop-golang-microservices/internal/pkg/gorm_pgsql"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/grpc"
	echoserver "github.com/meysamhadeli/shop-golang-microservices/internal/pkg/http/echo/server"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/logger"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/otel"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/rabbitmq"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var configPath string

type Config struct {
	ServiceName  string                        `mapstructure:"serviceName"`
	Logger       *logger.LoggerConfig          `mapstructure:"logger"`
	Rabbitmq     *rabbitmq.RabbitMQConfig      `mapstructure:"rabbitmq"`
	Echo         *echoserver.EchoConfig        `mapstructure:"echo"`
	Grpc         *grpc.GrpcConfig              `mapstructure:"grpc"`
	GormPostgres *gormpgsql.GormPostgresConfig `mapstructure:"gormPostgres"`
	Jaeger       *otel.JaegerConfig            `mapstructure:"jaeger"`
}

func init() {
	flag.StringVar(&configPath, "config", "", "products write microservice config path")
}

func InitConfig() (*Config, *logger.LoggerConfig, *otel.JaegerConfig, *gormpgsql.GormPostgresConfig,
	*grpc.GrpcConfig, *echoserver.EchoConfig, *rabbitmq.RabbitMQConfig, error) {

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	if configPath == "" {
		configPathFromEnv := os.Getenv("CONFIG_PATH")
		if configPathFromEnv != "" {
			configPath = configPathFromEnv
		} else {
			//https://stackoverflow.com/questions/31873396/is-it-possible-to-get-the-current-root-of-package-structure-as-a-string-in-golan
			//https://stackoverflow.com/questions/18537257/how-to-get-the-directory-of-the-currently-running-file
			d, err := dirname()
			if err != nil {
				return nil, nil, nil, nil, nil, nil, nil, err
			}

			configPath = d
		}
	}

	cfg := &Config{}

	viper.SetConfigName(fmt.Sprintf("config.%s", env))
	viper.AddConfigPath(configPath)
	viper.SetConfigType("json")

	if err := viper.ReadInConfig(); err != nil {
		return nil, nil, nil, nil, nil, nil, nil, errors.Wrap(err, "viper.ReadInConfig")
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, nil, nil, nil, nil, nil, nil, errors.Wrap(err, "viper.Unmarshal")
	}

	return cfg, cfg.Logger, cfg.Jaeger, cfg.GormPostgres, cfg.Grpc, cfg.Echo, cfg.Rabbitmq, nil
}

func GetMicroserviceName(serviceName string) string {
	return fmt.Sprintf("%s", strings.ToUpper(serviceName))
}

func filename() (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.New("unable to get the current filename")
	}
	return filename, nil
}

func dirname() (string, error) {
	filename, err := filename()
	if err != nil {
		return "", err
	}
	return filepath.Dir(filename), nil
}
