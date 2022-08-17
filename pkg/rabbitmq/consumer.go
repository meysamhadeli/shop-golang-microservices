package rabbitmq

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/meysamhadeli/shop-golang-microservices/pkg/logger"
	"github.com/streadway/amqp"
	"reflect"
)

type IConsumer interface {
	ConsumeMessage(msg interface{}) (*<-chan amqp.Delivery, error)
}

type consumer struct {
	cfg  *RabbitMQConfig
	conn *amqp.Connection
	log  logger.ILogger
}

func (p consumer) ConsumeMessage(msg interface{}) (*<-chan amqp.Delivery, error) {

	ch, err := p.conn.Channel()
	if err != nil {
		return nil, err
	}

	defer ch.Close()

	typeName := reflect.TypeOf(msg).Elem().Name()
	snakeTypeName := strcase.ToSnake(typeName)

	if err != nil {
		return nil, err
	}

	err = ch.ExchangeDeclare(
		snakeTypeName, // name
		p.cfg.Kind,    // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	)

	if err != nil {
		return nil, err
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
		return nil, err
	}

	err = ch.QueueBind(
		q.Name,        // queue name
		snakeTypeName, // routing key
		snakeTypeName, // exchange
		false,
		nil)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	return &msgs, nil
}

func NewConsumer(cfg *RabbitMQConfig, conn *amqp.Connection, log logger.ILogger) *consumer {
	return &consumer{cfg: cfg, conn: conn, log: log}
}
