package rabbitmq

import (
	"context"
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/meysamhadeli/shop-golang-microservices/pkg/logger"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
	"reflect"
)

type IConsumer interface {
	ConsumeMessage(msg interface{}) error
}

type consumer struct {
	cfg     *RabbitMQConfig
	conn    *amqp.Connection
	log     logger.ILogger
	handler func(queue string, msg amqp.Delivery, err error)
}

func (c consumer) ConsumeMessage(ctx context.Context, msg interface{}) error {

	ch, err := c.conn.Channel()
	if err != nil {
		c.log.Error("Error in opening channel to consume message")
	}

	defer ch.Close()

	typeName := reflect.TypeOf(msg).Name()
	snakeTypeName := strcase.ToSnake(typeName)

	err = ch.ExchangeDeclare(
		snakeTypeName, // name
		c.cfg.Kind,    // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	)

	if err != nil {
		c.log.Error("Error in declaring exchange to consume message")
	}

	q, err := ch.QueueDeclare(
		fmt.Sprintf("%s_%s", snakeTypeName, "queue"), // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)

	if err != nil {
		c.log.Error("Error in declaring queue to consume message")
	}

	err = ch.QueueBind(
		q.Name,        // queue name
		snakeTypeName, // routing key
		snakeTypeName, // exchange
		false,
		nil)
	if err != nil {
		c.log.Error("Error in binding queue to consume message")
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)

	if err != nil {
		c.log.Error("Error in consuming message")
	}

	for {
		select {
		case <-ctx.Done():
			c.log.Errorf("consumer ctx done: %v", ctx.Err())
			return ctx.Err()

		case m, ok := <-msgs:
			if !ok {
				c.log.Errorf("NOT OK deliveries channel closed for queue: %s", q.Name)
				return errors.New("deliveries channel closed")
			}

			c.handler(q.Name, m, nil)
		}
	}

	return nil
}

func NewConsumer(cfg *RabbitMQConfig, conn *amqp.Connection, log logger.ILogger, handler func(queue string, msg amqp.Delivery, err error)) *consumer {
	return &consumer{cfg: cfg, conn: conn, log: log, handler: handler}
}
