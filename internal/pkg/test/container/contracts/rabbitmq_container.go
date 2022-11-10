package contracts

import (
	"context"
	"github.com/streadway/amqp"
	"testing"
)

type RabbitMQContainerOptions struct {
	Host        string
	VirtualHost string
	Ports       []string
	HostPort    int
	UserName    string
	Password    string
	ImageName   string
	Name        string
	Tag         string
}

type RabbitMQContainer interface {
	Start(ctx context.Context, t *testing.T, options ...*RabbitMQContainerOptions) (*amqp.Connection, error, func())
	Cleanup(ctx context.Context) error
}
