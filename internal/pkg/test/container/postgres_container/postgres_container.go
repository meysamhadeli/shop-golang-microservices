package postgrescontainer

import (
	"context"
	"fmt"
	"github.com/docker/go-connections/nat"
	gormpgsql "github.com/meysamhadeli/shop-golang-microservices/internal/pkg/gorm_pgsql"
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

func Start(ctx context.Context) (*gorm.DB, *PostgresContainer, error) {

	defaultPostgresOptions, err := getDefaultPostgresTestContainers()
	if err != nil {
		return nil, nil, err
	}

	req := getContainerRequest(defaultPostgresOptions)

	postgresContainer, err := testcontainers.GenericContainer(
		ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		})

	go func() {
		for {
			select {
			case <-ctx.Done():
				postgresContainer.Terminate(ctx)
			}
		}
	}()

	if err != nil {
		return nil, nil, err
	}

	host, err := postgresContainer.Host(ctx)
	if err != nil {

		return nil, nil, fmt.Errorf("failed to get container host: %v", err)
	}

	realPort, err := postgresContainer.MappedPort(ctx, defaultPostgresOptions.Port)

	if err != nil {
		return nil, nil, fmt.Errorf("failed to get exposed container port: %v", err)
	}

	containerPort := realPort.Int()

	var gormConfig = &gormpgsql.GormPostgresConfig{
		Port:     containerPort,
		Host:     host,
		DBName:   defaultPostgresOptions.Database,
		User:     defaultPostgresOptions.UserName,
		Password: defaultPostgresOptions.Password,
		SSLMode:  false,
	}

	db, err := gormpgsql.NewGorm(gormConfig)

	return db, &PostgresContainer{Container: postgresContainer}, nil
}

func getContainerRequest(opts *PostgresContainerOptions) testcontainers.ContainerRequest {

	containerReq := testcontainers.ContainerRequest{
		Image:        fmt.Sprintf("%s:%s", opts.ImageName, opts.Tag),
		Hostname:     opts.Host,
		ExposedPorts: []string{string(opts.Port)},
		WaitingFor:   wait.ForListeningPort(opts.Port).WithStartupTimeout(opts.Timeout),
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
