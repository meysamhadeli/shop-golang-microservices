package rabbitmq

import (
	"context"
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
}

// Initialize new channel for rabbitmq
func NewRabbitMQConn(cfg *RabbitMQConfig, ctx context.Context) (*amqp.Connection, error) {

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
		return nil, err
	}

	go func() {
		select {
		case <-ctx.Done():
			defer func(conn *amqp.Connection) {
				err := conn.Close()
				if err != nil {
					log.Error("Failed to close connection")
				}
			}(conn)
			log.Info("Connection is closed")
		}
	}()

	log.Info("Connected to RabbitMQ")

	return conn, nil
}
