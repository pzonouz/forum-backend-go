package services_test

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"

	"forum-backend-go/ineternal/utils"
)

func TestUserService(t *testing.T) {
	t.Parallel()
	require := require.New(t)
	dsn := os.Getenv("dsn")

	if dsn == "" {
		dsn = "host=localhost port=5432 user=root password=secret dbname=forum sslmode=disable"
	}

	conn, err := sql.Open("pgx", dsn)
	require.NoError(err)
	err = conn.Ping()
	require.NoError(err)

	defer conn.Close()
	_, err = conn.Exec(utils.CreateUserTableQueryTest)
	require.NoError(err)
	_, err = conn.Exec(utils.DeleteUserTableQueryTest)
	require.NoError(err)
}
