package contracts

import (
	"context"
	"github.com/docker/go-connections/nat"
	"github.com/streadway/amqp"
	"testing"
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
	Timeout     time.Duration
}

type RabbitMQContainer interface {
	Start(ctx context.Context, t *testing.T, options ...*RabbitMQContainerOptions) (*amqp.Connection, error, func())
}
