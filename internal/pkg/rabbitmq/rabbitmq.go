package rabbitmq

import (
	"fmt"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/config_options"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

// Initialize new channel for rabbitmq
func NewRabbitMQConn(cfg *config_options.Config) (*amqp.Connection, error, func()) {

	connAddr := fmt.Sprintf(
		"amqp://%s:%s@%s:%d/",
		cfg.Rabbitmq.User,
		cfg.Rabbitmq.Password,
		cfg.Rabbitmq.Host,
		cfg.Rabbitmq.Port,
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
