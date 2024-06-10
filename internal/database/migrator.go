package database

// TODO:
// 1. Copy the "main schema" file to the migrations, so I can "migrate" my database from scratch to the latest version.
// 2. Implement a function that will run the migrations on a given database.

import (
	"context"
	"fmt"
	"log"
	"os/exec"

	_ "github.com/jackc/pgx/v5/stdlib" // Import pgx driver
	"github.com/pressly/goose/v3"
	"github.com/pro-posal/webserver/config"
)

func RunDockerContainer(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, "docker", "run", "-d", "--name", "TempDB",
		"-e", "POSTGRES_USER="+config.TestConfig.TestDatabase.UserName,
		"-e", "POSTGRES_PASSWORD="+config.TestConfig.TestDatabase.Password,
		"-e", "POSTGRES_DB="+config.TestConfig.TestDatabase.DbName,
		"-p", config.TestConfig.TestDatabase.Port+":5432", "postgres")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to start Docker container: %w\n%s", err, output)
	}

	containerID := string(output)
	return containerID, nil
}

func StopDockerContainer(ctx context.Context) {
	stopCmd := exec.CommandContext(ctx, "docker", "stop", "TempDB")
	err := stopCmd.Run()
	if err != nil {
		log.Fatalf("Failed to stop Docker container: %v", err)
	}

	// Remove the Docker container
	removeCmd := exec.CommandContext(ctx, "docker", "rm", "TempDB")
	err = removeCmd.Run()
	if err != nil {
		log.Fatalf("Failed to remove Docker container: %v", err)
	}
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
