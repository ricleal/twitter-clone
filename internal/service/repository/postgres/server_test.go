package postgres

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	pgMigrate "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type PostgresTestSuite struct {
	suite.Suite
	container *postgres.PostgresContainer
	dbURL     string
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run.
func TestPostgresTestSuite(t *testing.T) {
	suite.Run(t, new(PostgresTestSuite))
}

func (ts *PostgresTestSuite) SetupTest() {

	dbname := os.Getenv("DB_NAME") + "_test"
	user := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")

	ctx := context.Background()
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
		ts.T().Fatal(err)
	}
	ts.container = container

	mappedPort, err := container.MappedPort(ctx, "5432")
	if err != nil {
		ts.T().Fatal(err)
	}

	hostIP, err := container.Host(ctx)
	if err != nil {
		ts.T().Fatal(err)
	}

	uri := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, hostIP, mappedPort.Port(), dbname)
	ts.dbURL = uri
	os.Setenv("DATABASE_URL", uri)

	// run migrations
	s := New()

	// db, err := sql.Open("postgres", "postgres://localhost:5432/database?sslmode=enable")
	// driver, err := postgres.WithInstance(db, &postgres.Config{})
	driver, err := pgMigrate.WithInstance(s.dbConn, &pgMigrate.Config{})
	if err != nil {
		ts.T().Fatal(fmt.Errorf("failed to create migration driver: %w", err))
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://../../../../migrations",
		"postgres", driver)
	if err != nil {
		ts.T().Fatal(fmt.Errorf("failed to create migration instance: %w", err))
	}
	m.Up() // or m.Step(2) if you want to explicitly set the number of migrations to run
}

func (ts *PostgresTestSuite) TestPostgres() {
	s := New()
	require.NotNil(ts.T(), s.dbConn)
	err := s.dbConn.Ping()
	require.NoError(ts.T(), err)

	// check the existing tables
	rows, err := s.dbConn.Query("SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'")
	require.NoError(ts.T(), err)
	defer rows.Close()
	var tables []string
	for rows.Next() {
		var table string
		err := rows.Scan(&table)
		require.NoError(ts.T(), err)
		tables = append(tables, table)
	}
	require.Contains(ts.T(), tables, "schema_migrations")
	require.Contains(ts.T(), tables, "users")
	require.Contains(ts.T(), tables, "tweets")
	s.Close()

}

func (ts *PostgresTestSuite) TearDownTest() {
	ts.container.Terminate(context.Background())
}
