package postgrescontainer

import (
	"context"
	"fmt"
	"github.com/cenkalti/backoff/v4"
	"github.com/docker/go-connections/nat"
	gormpgsql "github.com/meysamhadeli/shop-golang-microservices/internal/pkg/gorm_pgsql"
	"github.com/pkg/errors"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/gorm"
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

type PostgresContainer struct {
	Container testcontainers.Container
}

func (c *PostgresContainer) Terminate(ctx context.Context) {
	if c.Container != nil {
		c.Container.Terminate(ctx)
	}
}

func Start(ctx context.Context) (*gorm.DB, *gormpgsql.GormPostgresConfig, *PostgresContainer, error) {

	defaultPostgresOptions, err := getDefaultPostgresTestContainers()
	if err != nil {
		return nil, nil, nil, err
	}

	req := getContainerRequest(defaultPostgresOptions)

	postgresContainer, err := testcontainers.GenericContainer(
		ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		})

	if err != nil {
		return nil, nil, nil, err
	}

	var gormDB *gorm.DB
	var gormConfig *gormpgsql.GormPostgresConfig

	bo := backoff.NewExponentialBackOff()
	bo.MaxElapsedTime = 10 * time.Second // Maximum time to retry
	maxRetries := 5                      // Number of retries (including the initial attempt)

	err = backoff.Retry(func() error {

		host, err := postgresContainer.Host(ctx)
		if err != nil {

			return errors.Errorf("failed to get container host: %v", err)
		}

		realPort, err := postgresContainer.MappedPort(ctx, "5432")

		if err != nil {
			return errors.Errorf("failed to get exposed container port: %v", err)
		}

		containerPort := realPort.Int()

		gormConfig = &gormpgsql.GormPostgresConfig{
			Port:     containerPort,
			Host:     host,
			DBName:   defaultPostgresOptions.Database,
			User:     defaultPostgresOptions.UserName,
			Password: defaultPostgresOptions.Password,
			SSLMode:  false,
		}
		gormDB, err = gormpgsql.NewGorm(gormConfig)
		if err != nil {
			return err
		}
		return nil
	}, backoff.WithMaxRetries(bo, uint64(maxRetries-1)))

	if err != nil {
		return nil, nil, nil, errors.Errorf("failed to create connection for postgres after retries: %v", err)
	}

	return gormDB, gormConfig, &PostgresContainer{Container: postgresContainer}, nil
}

func getContainerRequest(opts *PostgresContainerOptions) testcontainers.ContainerRequest {

	containerReq := testcontainers.ContainerRequest{
		Image:        fmt.Sprintf("%s:%s", opts.ImageName, opts.Tag),
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp"),
		Env: map[string]string{
			"POSTGRES_DB":       opts.Database,
			"POSTGRES_PASSWORD": opts.Password,
			"POSTGRES_USER":     opts.UserName,
		},
	}

	return containerReq
}

func getDefaultPostgresTestContainers() (*PostgresContainerOptions, error) {
	port, err := nat.NewPort("", "5432")
	if err != nil {
		return nil, fmt.Errorf("failed to build port: %v", err)
	}

	return &PostgresContainerOptions{
		Database:  "test_db",
		Port:      port,
		Host:      "localhost",
		UserName:  "testcontainers",
		Password:  "testcontainers",
		Tag:       "latest",
		ImageName: "postgres",
		Name:      "postgresql-testcontainer",
		Timeout:   5 * time.Minute,
	}, nil
}
