package contracts

import (
	"context"
	"github.com/docker/go-connections/nat"
	"gorm.io/gorm"
	"testing"
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

type GormContainer interface {
	Start(ctx context.Context, t *testing.T, options ...*PostgresContainerOptions) (*gorm.DB, error)
}
