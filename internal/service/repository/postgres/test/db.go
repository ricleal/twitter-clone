package test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/golang-migrate/migrate/v4"
	pgMigrate "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupContainer(ctx context.Context) (*postgres.PostgresContainer, error) {

	dbname := os.Getenv("DB_NAME") + "_test"
	user := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	port := os.Getenv("DB_PORT")

	container, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("docker.io/postgres:15.3"),
		postgres.WithDatabase(dbname),
		postgres.WithUsername(user),
		postgres.WithPassword(password),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").WithOccurrence(2).WithStartupTimeout(5*time.Second),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to run postgres container: %w", err)
	}

	mappedPort, err := container.MappedPort(ctx, nat.Port(port))
	if err != nil {
		return nil, fmt.Errorf("failed to get mapped port: %w", err)
	}

	hostIP, err := container.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get host IP: %w", err)
	}

	uri := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, hostIP, mappedPort.Port(), dbname)
	os.Setenv("DB_URL", uri)
	return container, nil
}

func setupMigrations(ctx context.Context) error {
	dbURL := os.Getenv("DB_URL")

	if dbURL == "" {
		return fmt.Errorf("DB_URL not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	driver, err := pgMigrate.WithInstance(db, &pgMigrate.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	migrationsPath := os.Getenv("MIGRATIONS_PATH")
	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}
	m.Up()
	err = db.Close()
	if err != nil {
		return fmt.Errorf("failed to close database: %w", err)
	}
	return nil
}

func SetupDB(ctx context.Context) (*postgres.PostgresContainer, error) {
	container, err := setupContainer(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to setup container: %w", err)
	}
	err = setupMigrations(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to setup migrations: %w", err)
	}
	return container, nil
}

func TeardownDB(ctx context.Context, container *postgres.PostgresContainer) error {
	err := container.Terminate(ctx)
	if err != nil {
		return fmt.Errorf("failed to terminate container: %w", err)
	}
	return nil
}
