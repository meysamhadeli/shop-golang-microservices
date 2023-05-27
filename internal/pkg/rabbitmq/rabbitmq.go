package rabbitmq

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"time"
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
		log.Errorf("Failed to connect to RabbitMQ: %v. Connection information: %s", err, connAddr)
		return nil, err
	}

	log.Info("Connected to RabbitMQ")

	go func() {
		select {
		case <-ctx.Done():
			err := conn.Close()
			if err != nil {
				log.Error("Failed to close RabbitMQ connection")
			}
			log.Info("RabbitMQ connection is closed")
		}
	}()

	return conn, err
}

func retryWithBackoff(fn func() error, retryInterval time.Duration, maxRetries int) error {
	for i := 0; i < maxRetries; i++ {
		err := fn()
		if err == nil {
			return nil
		}

		time.Sleep(retryInterval)
		retryInterval *= 2 // exponential backoff
	}

	return fmt.Errorf("maximum number of retries reached")
}
