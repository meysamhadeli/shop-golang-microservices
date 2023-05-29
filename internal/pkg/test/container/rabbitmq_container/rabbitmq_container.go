package rabbitmqcontainer

import (
	"context"
	"fmt"
	"github.com/cenkalti/backoff/v4"
	"github.com/docker/go-connections/nat"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/rabbitmq"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"time"
)

type RabbitMQContainerOptions struct {
	Host        string
	VirtualHost string
	Port        nat.Port
	UserName    string
	Password    string
	ImageName   string
	Name        string
	Tag         string
	Exchange    string
	Timeout     time.Duration
	Kind        string
}

type RabbitmqContainer struct {
	Container testcontainers.Container
}

func (c *RabbitmqContainer) Terminate(ctx context.Context) {
	if c.Container != nil {
		c.Container.Terminate(ctx)
	}
}

// ref: https://github.com/romnn/testcontainers/blob/60ec1eb7563985ae83e51bb04ca3c67236787a26/rabbitmq/rabbitmq.go
func Start(ctx context.Context) (*amqp.Connection, *rabbitmq.RabbitMQConfig, *RabbitmqContainer, error) {

	defaultRabbitmqOptions, err := getDefaultRabbitMQTestContainers()
	if err != nil {
		return nil, nil, nil, err
	}

	req := getContainerRequest(defaultRabbitmqOptions)

	rmqContainer, err := testcontainers.GenericContainer(
		ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		})

	if err != nil {
		return nil, nil, nil, err
	}

	var conn *amqp.Connection
	var rabbitmqConfig *rabbitmq.RabbitMQConfig

	bo := backoff.NewExponentialBackOff()
	bo.MaxElapsedTime = 10 * time.Second // Maximum time to retry
	maxRetries := 5                      // Number of retries (including the initial attempt)

	err = backoff.Retry(func() error {
		host, err := rmqContainer.Host(ctx)
		if err != nil {
			return errors.Errorf("failed to get container host: %v", err)
		}

		realPort, err := rmqContainer.MappedPort(ctx, defaultRabbitmqOptions.Port)

		if err != nil {
			return errors.Errorf("failed to get exposed container port: %v", err)
		}

		log.Info(realPort)
		containerPort := realPort.Int()

		rabbitmqConfig = &rabbitmq.RabbitMQConfig{
			User:         defaultRabbitmqOptions.UserName,
			Password:     defaultRabbitmqOptions.Password,
			Host:         host,
			Port:         containerPort,
			ExchangeName: defaultRabbitmqOptions.Exchange,
			Kind:         defaultRabbitmqOptions.Kind,
		}

		conn, err = rabbitmq.NewRabbitMQConn(rabbitmqConfig, ctx)
		if err != nil {
			log.Errorf("Failed to create connection for rabbitmq: %v", err)
			return err
		}
		return nil
	}, backoff.WithMaxRetries(bo, uint64(maxRetries-1)))

	if err != nil {
		return nil, nil, nil, errors.Errorf("failed to create connection for rabbitmq after retries: %v", err)
	}

	return conn, rabbitmqConfig, &RabbitmqContainer{Container: rmqContainer}, err
}

func getContainerRequest(opts *RabbitMQContainerOptions) testcontainers.ContainerRequest {

	containerReq := testcontainers.ContainerRequest{
		Image:        fmt.Sprintf("%s:%s", opts.ImageName, opts.Tag),
		ExposedPorts: []string{"5672/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections"),
		Env: map[string]string{
			"RABBITMQ_DEFAULT_USER": opts.UserName,
			"RABBITMQ_DEFAULT_PASS": opts.Password,
		},
	}

	return containerReq
}

func getDefaultRabbitMQTestContainers() (*RabbitMQContainerOptions, error) {
	port, err := nat.NewPort("", "5672")
	if err != nil {
		return nil, fmt.Errorf("failed to build port: %v", err)
	}

	return &RabbitMQContainerOptions{
		Port:        port,
		Host:        "localhost",
		VirtualHost: "/",
		UserName:    "guest",
		Password:    "guest",
		Tag:         "3-management",
		ImageName:   "rabbitmq",
		Exchange:    "test",
		Kind:        "topic",
		Name:        "rabbitmq-testcontainers",
		Timeout:     5 * time.Minute,
	}, nil
}
