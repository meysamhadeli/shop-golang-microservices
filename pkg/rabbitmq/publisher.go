package rabbitmq

import (
	"github.com/iancoleman/strcase"
	jsoniter "github.com/json-iterator/go"
	"github.com/meysamhadeli/shop-golang-microservices/pkg/logger"
	uuid "github.com/satori/go.uuid"
	"github.com/streadway/amqp"
	"reflect"
	"time"
)

type IPublisher interface {
	PublishMessage(msg interface{}) error
}

type publisher struct {
	cfg  *RabbitMQConfig
	conn *amqp.Connection
	log  logger.ILogger
}

func (p publisher) PublishMessage(msg interface{}) error {

	ch, err := p.conn.Channel()
	if err != nil {
		p.log.Error("Error in opening channel to publish message")
	}

	defer ch.Close()

	data, err := jsoniter.Marshal(msg)
	typeName := reflect.TypeOf(msg).Elem().Name()
	snakeTypeName := strcase.ToSnake(typeName)

	if err != nil {
		p.log.Error("Error in marshalling message to publish message")
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
		p.log.Error("Error in declaring exchange to publish message")
	}

	publishingMsg := amqp.Publishing{
		Body:         data,
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		MessageId:    uuid.NewV4().String(),
		Timestamp:    time.Now(),
	}

	err = ch.Publish(snakeTypeName, snakeTypeName, false, false, publishingMsg)

	if err != nil {
		p.log.Error("Error in publishing message")
	}

	p.log.Infof("Published message: %s", publishingMsg)
	return nil
}

func NewPublisher(cfg *RabbitMQConfig, conn *amqp.Connection, log logger.ILogger) *publisher {
	return &publisher{cfg: cfg, conn: conn, log: log}
}
