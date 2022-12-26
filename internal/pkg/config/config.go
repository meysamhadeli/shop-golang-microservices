package config

import (
	"flag"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "", "products write microservice config path")
}

type LoggerConfig struct {
	LogLevel string `mapstructure:"level"`
}

type RabbitMQConfig struct {
	Host         string
	Port         int
	User         string
	Password     string
	ExchangeName string
	Kind         string
}

type EchoConfig struct {
	Port                string   `mapstructure:"port" validate:"required"`
	Development         bool     `mapstructure:"development"`
	BasePath            string   `mapstructure:"basePath" validate:"required"`
	DebugErrorsResponse bool     `mapstructure:"debugErrorsResponse"`
	IgnoreLogUrls       []string `mapstructure:"ignoreLogUrls"`
	Timeout             int      `mapstructure:"timeout"`
	Host                string   `mapstructure:"host"`
}

type GrpcConfig struct {
	Port        string `mapstructure:"port"`
	Host        string `mapstructure:"host"`
	Development bool   `mapstructure:"development"`
}

type GormPostgresConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	DBName   string `mapstructure:"dbName"`
	SSLMode  bool   `mapstructure:"sslMode"`
	Password string `mapstructure:"password"`
}

type JaegerConfig struct {
	Server      string `mapstructure:"server"`
	ServiceName string `mapstructure:"serviceName"`
	TracerName  string `mapstructure:"tracerName"`
}

type Config struct {
	ServiceName  string              `mapstructure:"serviceName"`
	Logger       *LoggerConfig       `mapstructure:"logger"`
	Rabbitmq     *RabbitMQConfig     `mapstructure:"rabbitmq"`
	Echo         *EchoConfig         `mapstructure:"echo"`
	Grpc         *GrpcConfig         `mapstructure:"grpc"`
	Context      Context             `mapstructure:"context"`
	GormPostgres *GormPostgresConfig `mapstructure:"gormPostgres"`
	Jaeger       *JaegerConfig       `mapstructure:"jaeger"`
}

type Context struct {
	Timeout int `mapstructure:"timeout"`
}

func InitConfig() (*Config, error) {

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
				return nil, err
			}

			configPath = d
		}
	}

	cfg := &Config{}

	viper.SetConfigName(fmt.Sprintf("config.%s", env))
	viper.AddConfigPath(configPath)
	viper.SetConfigType("json")

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
