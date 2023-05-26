package rabbitmqcontainer

import (
	"context"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/rabbitmq"
	log "github.com/sirupsen/logrus"
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
func Start(ctx context.Context) (*rabbitmq.RabbitMQConfig, *RabbitmqContainer, error) {

	defaultRabbitmqOptions, err := getDefaultRabbitMQTestContainers()
	if err != nil {
		return nil, nil, err
	}

	req := getContainerRequest(defaultRabbitmqOptions)

	rmqContainer, err := testcontainers.GenericContainer(
		ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		})

	if err != nil {
		return nil, nil, err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				rmqContainer.Terminate(ctx)
			}
		}
	}()

	host, err := rmqContainer.Host(ctx)
	if err != nil {

		return nil, nil, fmt.Errorf("failed to get container host: %v", err)
	}

	realPort, err := rmqContainer.MappedPort(ctx, defaultRabbitmqOptions.Port)

	if err != nil {
		return nil, nil, fmt.Errorf("failed to get exposed container port: %v", err)
	}

	log.Info(realPort)
	containerPort := realPort.Int()

	var rabbitmqConfig = &rabbitmq.RabbitMQConfig{
		User:         defaultRabbitmqOptions.UserName,
		Password:     defaultRabbitmqOptions.Password,
		Host:         host,
		Port:         containerPort,
		ExchangeName: defaultRabbitmqOptions.Exchange,
		Kind:         defaultRabbitmqOptions.Kind,
	}

	return rabbitmqConfig, &RabbitmqContainer{Container: rmqContainer}, nil
}

func getContainerRequest(opts *RabbitMQContainerOptions) testcontainers.ContainerRequest {

	containerReq := testcontainers.ContainerRequest{
		Image:        fmt.Sprintf("%s:%s", opts.ImageName, opts.Tag),
		Hostname:     opts.Host,
		ExposedPorts: []string{string(opts.Port)},
		WaitingFor:   wait.ForListeningPort(opts.Port).WithStartupTimeout(opts.Timeout),
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
