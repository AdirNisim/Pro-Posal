package database

import (
	"context"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib" // Import pgx driver
	"github.com/pressly/goose/v3"
	"github.com/pro-posal/webserver/config"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func RunDockerContainer(ctx context.Context) (string, string, testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		Image:        "postgres:latest", // or the specific image you need
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_DB":       config.TestConfig.TestDatabase.DbName,
			"POSTGRES_USER":     config.TestConfig.TestDatabase.UserName,
			"POSTGRES_PASSWORD": config.TestConfig.TestDatabase.Password,
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}
	postgres, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return "", "", nil, err
	}

	host, err := postgres.Host(ctx)
	if err != nil {
		return "", "", postgres, err
	}

	mappedPort, err := postgres.MappedPort(ctx, "5432")
	if err != nil {
		return "", "", postgres, err
	}

	log.Printf("PostgreSQL container started on host: %s, port: %s", host, mappedPort.Port())
	return host, mappedPort.Port(), postgres, nil
}

func RunMigrations() error {
	dbString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.TestConfig.TestDatabase.UserName,
		config.TestConfig.TestDatabase.Password,
		config.TestConfig.TestDatabase.Host,
		config.TestConfig.TestDatabase.Port,
		config.TestConfig.TestDatabase.DbName)

	db, err := goose.OpenDBWithDriver("pgx", dbString)
	if err != nil {
		return fmt.Errorf("failed to open DB: %w", err)
	}
	defer db.Close()

	err = goose.Up(db, "../../migrations")
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}
