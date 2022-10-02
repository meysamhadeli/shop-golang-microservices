package rabbitmq

import (
	"context"
	"fmt"
	"github.com/iancoleman/strcase"
	jsoniter "github.com/json-iterator/go"
	"github.com/meysamhadeli/shop-golang-microservices/pkg/logger"
	open_telemetry "github.com/meysamhadeli/shop-golang-microservices/pkg/open-telemetry"
	"github.com/streadway/amqp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"reflect"
	"runtime"
	"strings"
	"time"
)

type IConsumer interface {
	ConsumeMessage(ctx context.Context, msg interface{}) error
}

type consumer struct {
	cfg          *RabbitMQConfig
	conn         *amqp.Connection
	log          logger.ILogger
	handler      func(queue string, msg amqp.Delivery) error
	jaegerTracer trace.Tracer
}

func (c consumer) ConsumeMessage(ctx context.Context, msg interface{}) (error, func()) {

	strName := strings.Split(runtime.FuncForPC(reflect.ValueOf(c.handler).Pointer()).Name(), ".")
	var consumerHandlerName = strName[len(strName)-1]

	ch, err := c.conn.Channel()
	if err != nil {
		c.log.Error("Error in opening channel to consume message")
		return err, nil
	}

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
		return err, nil
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
		return err, nil
	}

	err = ch.QueueBind(
		q.Name,        // queue name
		snakeTypeName, // routing key
		snakeTypeName, // exchange
		false,
		nil)
	if err != nil {
		c.log.Error("Error in binding queue to consume message")
		return err, nil
	}

	deliveries, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)

	if err != nil {
		c.log.Error("Error in consuming message")
		return err, nil
	}

	forever := make(chan error)

	go func() {
		for delivery := range deliveries {

			// Extract headers
			ctx = open_telemetry.ExtractAMQPHeaders(ctx, delivery.Headers)

			err := c.handler(q.Name, delivery)
			if err != nil {
				c.log.Error(err.Error())
			}

			_, span := c.jaegerTracer.Start(ctx, consumerHandlerName)

			h, err := jsoniter.Marshal(delivery.Headers)

			if err != nil {
				c.log.Errorf("Error in marshalling headers in consumer: %v", string(h))
			}

			span.SetAttributes(attribute.Key("message-id").String(delivery.MessageId))
			span.SetAttributes(attribute.Key("correlation-id").String(delivery.CorrelationId))
			span.SetAttributes(attribute.Key("queue").String(q.Name))
			span.SetAttributes(attribute.Key("exchange").String(delivery.Exchange))
			span.SetAttributes(attribute.Key("routing-key").String(delivery.RoutingKey))
			span.SetAttributes(attribute.Key("ack").Bool(true))
			span.SetAttributes(attribute.Key("timestamp").String(delivery.Timestamp.String()))
			span.SetAttributes(attribute.Key("body").String(string(delivery.Body)))
			span.SetAttributes(attribute.Key("headers").String(string(h)))

			// Cannot use defer inside a for loop
			time.Sleep(1 * time.Millisecond)
			span.End()

			err = delivery.Ack(false)
			if err != nil {
				c.log.Errorf("We didn't get a ack for delivery: %v", string(delivery.Body))
			}
		}
	}()

	c.log.Info("Waiting for messages. To exit press CTRL+C")

	forever <- nil

	return nil, func() {
		_ = ch.Close()
	}
}

func NewConsumer(cfg *RabbitMQConfig, conn *amqp.Connection, log logger.ILogger, jaegerTracer trace.Tracer, handler func(queue string, msg amqp.Delivery) error) *consumer {
	return &consumer{cfg: cfg, conn: conn, log: log, jaegerTracer: jaegerTracer, handler: handler}
}
