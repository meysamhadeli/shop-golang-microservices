package rabbitmq

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/meysamhadeli/shop-golang-microservices/pkg/logger"
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

func (c consumer) ConsumeMessage(msg interface{}) error {

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

	forever := make(chan struct{})

	go func() {
		for m := range msgs {
			c.handler(q.Name, m, nil)
		}
	}()

	c.log.Info("Waiting for messages. To exit press CTRL+C")

	<-forever

	return nil
}

func NewConsumer(cfg *RabbitMQConfig, conn *amqp.Connection, log logger.ILogger, handler func(queue string, msg amqp.Delivery, err error)) *consumer {
	return &consumer{cfg: cfg, conn: conn, log: log, handler: handler}
}
