//go:build integration
// +build integration

package postgres_test

import (
	"context"
	"testing"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	testcontainers "github.com/testcontainers/testcontainers-go/modules/postgres"

	"github.com/ricleal/twitter-clone/internal/service/repository/postgres"
	"github.com/ricleal/twitter-clone/internal/service/repository/postgres/test"
)

type PostgresTestSuite struct {
	suite.Suite
	container *testcontainers.PostgresContainer
	ctx       context.Context
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run.
func TestPostgresTestSuite(t *testing.T) {
	suite.Run(t, new(PostgresTestSuite))
}

func (ts *PostgresTestSuite) SetupTest() {
	var err error
	ts.ctx = context.Background()
	ts.container, err = test.SetupDB(ts.ctx)
	require.NoError(ts.T(), err)
}

func (ts *PostgresTestSuite) TearDownTest() {
	err := test.TeardownDB(ts.ctx, ts.container)
	require.NoError(ts.T(), err)
}

func (ts *PostgresTestSuite) TestPostgres() {
	s, err := postgres.NewHandler(ts.ctx)
	require.NoError(ts.T(), err)
	require.NotNil(ts.T(), s.DB())
	err = s.DB().Ping()
	require.NoError(ts.T(), err)

	// check the existing tables
	rows, err := s.DB().Query("SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'")
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
