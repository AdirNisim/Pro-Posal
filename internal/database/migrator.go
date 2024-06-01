package database

// TODO:
// 1. Copy the "main schema" file to the migrations, so I can "migrate" my database from scratch to the latest version.
// 2. Implement a function that will run the migrations on a given database.

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" // Import pgx driver
	"github.com/pressly/goose/v3"
)

const (
	dbUser     = "user"
	dbPassword = "123456"
	dbName     = "testdb"
	dbHost     = "localhost"
	dbPort     = "5543"
)

func TestMain() {
	ctx := context.Background()

	// Step 1: Run a new Docker container
	_, err := runDockerContainer(ctx)
	if err != nil {
		log.Fatalf("Failed to run Docker container: %v", err)
	}
	defer stopDockerContainer(ctx)

	// Step 2: Wait for the database to be ready
	time.Sleep(2 * time.Second)

	// Step 3: Run migrations
	if err := runMigrations(); err != nil {
		fmt.Printf("Failed to run migrations: %v", err)
		return
	}

	fmt.Println("Database setup completed successfully")
}

func runDockerContainer(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, "docker", "run", "-d", "--name", "TempDB",
		"-e", "POSTGRES_USER="+dbUser,
		"-e", "POSTGRES_PASSWORD="+dbPassword,
		"-e", "POSTGRES_DB="+dbName,
		"-p", dbPort+":5432", "postgres")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to start Docker container: %w\n%s", err, output)
	}

	containerID := string(output)
	return containerID, nil
}

func stopDockerContainer(ctx context.Context) {
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

func runMigrations() error {
	dbString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)
	db, err := goose.OpenDBWithDriver("pgx", dbString)
	if err != nil {
		return fmt.Errorf("failed to open DB: %w", err)
	}
	defer db.Close()

	err = goose.Up(db, "../migrations")
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}
