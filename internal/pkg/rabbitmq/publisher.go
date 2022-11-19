package rabbitmq

import (
	"context"
	"github.com/ahmetb/go-linq/v3"
	"github.com/iancoleman/strcase"
	jsoniter "github.com/json-iterator/go"
	"github.com/labstack/echo/v4"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/logger"
	open_telemetry "github.com/meysamhadeli/shop-golang-microservices/internal/pkg/open-telemetry"
	uuid "github.com/satori/go.uuid"
	"github.com/streadway/amqp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"reflect"
	"time"
)

type IPublisher interface {
	PublishMessage(ctx context.Context, msg interface{}) error
	IsPublished(msg interface{}) bool
}

var publishedMessages []string

type publisher struct {
	cfg          *RabbitMQConfig
	conn         *amqp.Connection
	log          logger.ILogger
	jaegerTracer trace.Tracer
}

func (p publisher) PublishMessage(ctx context.Context, msg interface{}) error {

	data, err := jsoniter.Marshal(msg)

	if err != nil {
		p.log.Error("Error in marshalling message to publish message")
		return err
	}

	typeName := reflect.TypeOf(msg).Elem().Name()
	snakeTypeName := strcase.ToSnake(typeName)

	ctx, span := p.jaegerTracer.Start(ctx, typeName)
	defer span.End()

	// Inject the context in the headers
	headers := open_telemetry.InjectAMQPHeaders(ctx)

	channel, err := p.conn.Channel()
	if err != nil {
		p.log.Error("Error in opening channel to consume message")
		return err
	}

	defer channel.Close()

	err = channel.ExchangeDeclare(
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
		return err
	}

	correlationId := ""

	if ctx.Value(echo.HeaderXCorrelationID) != nil {
		correlationId = ctx.Value(echo.HeaderXCorrelationID).(string)
	}

	publishingMsg := amqp.Publishing{
		Body:          data,
		ContentType:   "application/json",
		DeliveryMode:  amqp.Persistent,
		MessageId:     uuid.NewV4().String(),
		Timestamp:     time.Now(),
		CorrelationId: correlationId,
		Headers:       headers,
	}

	err = channel.Publish(snakeTypeName, snakeTypeName, false, false, publishingMsg)

	if err != nil {
		p.log.Error("Error in publishing message")
		return err
	}

	publishedMessages = append(publishedMessages, snakeTypeName)

	h, err := jsoniter.Marshal(headers)

	if err != nil {
		p.log.Error("Error in marshalling headers to publish message")
		return err
	}

	p.log.Infof("Published message: %s", publishingMsg.Body)
	span.SetAttributes(attribute.Key("message-id").String(publishingMsg.MessageId))
	span.SetAttributes(attribute.Key("correlation-id").String(publishingMsg.CorrelationId))
	span.SetAttributes(attribute.Key("exchange").String(snakeTypeName))
	span.SetAttributes(attribute.Key("kind").String(p.cfg.Kind))
	span.SetAttributes(attribute.Key("content-type").String("application/json"))
	span.SetAttributes(attribute.Key("timestamp").String(publishingMsg.Timestamp.String()))
	span.SetAttributes(attribute.Key("body").String(string(publishingMsg.Body)))
	span.SetAttributes(attribute.Key("headers").String(string(h)))

	return nil
}

func (p publisher) IsPublished(msg interface{}) bool {

	typeName := reflect.TypeOf(msg).Name()
	snakeTypeName := strcase.ToSnake(typeName)
	isPublished := linq.From(publishedMessages).Contains(snakeTypeName)

	return isPublished
}

func NewPublisher(cfg *RabbitMQConfig, conn *amqp.Connection, log logger.ILogger, jaegerTracer trace.Tracer) *publisher {
	return &publisher{cfg: cfg, conn: conn, log: log, jaegerTracer: jaegerTracer}
}
