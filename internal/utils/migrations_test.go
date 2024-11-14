package utils_test

import (
	"context"
	"server/internal/storage/dbstorage"
	"server/internal/utils"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
)

var testDB *pgxpool.Pool

func TestMain(m *testing.M) {
	var err error
	testDB, err = dbstorage.Connection("postgres://postgres:root@localhost:5431/testdb?sslmode=disable")
	if err != nil {
		panic("failed to connect to test database")
	}
	defer testDB.Close()

	m.Run()
}

func TestMigrateUp(t *testing.T) {
	migrationPath := "./"
	err := utils.Migrate(testDB, migrationPath, "up")
	assert.NoError(t, err)

	row := testDB.QueryRow(context.Background(), "SELECT count(*) FROM metrics")
	var count int
	err = row.Scan(&count)
	assert.NoError(t, err)
	assert.Greater(t, count, 0)
}

func TestMigrateDown(t *testing.T) {
	migrationPath := "./"
	err := utils.Migrate(testDB, migrationPath, "down")
	assert.NoError(t, err)

	row := testDB.QueryRow(context.Background(), "SELECT to_regclass('metrics')")
	var tableName string
	err = row.Scan(&tableName)
	assert.NoError(t, err)
	assert.Equal(t, "metrics", tableName)
}
