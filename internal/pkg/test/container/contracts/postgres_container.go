package contracts

import (
	"github.com/docker/go-connections/nat"
	"time"
)

type PostgresContainerOptions struct {
	Database  string
	Host      string
	Port      nat.Port
	HostPort  int
	UserName  string
	Password  string
	ImageName string
	Name      string
	Tag       string
	Timeout   time.Duration
}
