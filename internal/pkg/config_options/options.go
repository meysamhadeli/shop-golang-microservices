package config_options

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
