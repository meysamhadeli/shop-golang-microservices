package contracts

import (
	"github.com/docker/go-connections/nat"
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
