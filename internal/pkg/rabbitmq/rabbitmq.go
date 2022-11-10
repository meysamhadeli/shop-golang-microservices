package rabbitmq

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type RabbitMQConfig struct {
	Host         string
	Port         int
	User         string
	Password     string
	ExchangeName string
	Kind         string
	VirtualHost  string `mapstructure:"virtualHost"`
}

// Initialize new channel for rabbitmq
func NewRabbitMQConn(cfg *RabbitMQConfig) (*amqp.Connection, error, func()) {

	connAddr := fmt.Sprintf(
		"amqp://%s:%s@%s:%d/",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
	)

	conn, err := amqp.Dial(connAddr)
	if err != nil {
		log.Error(err, "Failed to connect to RabbitMQ")
		return nil, err, nil
	}

	log.Info("Connected to RabbitMQ")

	return conn, nil, func() {
		_ = conn.Close()
	}
}
