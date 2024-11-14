package dbstorage

import (
	"context"
	"testing"

	"server/internal/dto"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

var testDB *pgxpool.Pool

func TestMain(m *testing.M) {
	var err error
	testDB, err = Connection("postgres://postgres:root@localhost:5431/testdb?sslmode=disable")
	if err != nil {
		panic("failed to connect to test database")
	}
	defer testDB.Close()

	m.Run()
}

func TestConnection(t *testing.T) {
	db, err := Connection("postgres://postgres:root@localhost:5431/testdb?sslmode=disable")
	require.NoError(t, err)
	assert.NotNil(t, db)
	db.Close()
}

func TestConnectionFall(t *testing.T) {
	db, err := Connection("postgres://postgres:root@uncorrecthost:5431/testdb?sslmode=disable")
	require.Error(t, err)
	assert.Nil(t, db)
}

func TestSetMetric(t *testing.T) {
	cleanMetricsDatabase(t, testDB)

	logger := zap.NewExample()
	defer logger.Sync()
	storage := New(testDB, logger)

	var value float64 = 123
	metric := dto.Metric{
		Name:  "test_metric",
		Type:  "gauge",
		Value: &value,
	}

	storedMetric, err := storage.SetMetric(metric)
	require.NoError(t, err)
	assert.NotNil(t, storedMetric, "metric should be stored successfully")
	assert.Equal(t, "test_metric", storedMetric.Name)
	assert.Equal(t, "gauge", storedMetric.Type)
	assert.Equal(t, value, *storedMetric.Value)
}

func TestSetUpdates(t *testing.T) {
	cleanMetricsDatabase(t, testDB)

	logger := zap.NewExample()
	defer logger.Sync()
	storage := New(testDB, logger)

	var delta int64 = 5
	metrics := []dto.Metric{
		{Name: "test_counter", Type: "counter", Value: new(float64), Delta: &delta},
	}

	storedMetrics, err := storage.SetUpdates(metrics)
	require.NoError(t, err)
	require.NotEmpty(t, storedMetrics, "metrics should be updated successfully")
	assert.Equal(t, "test_counter", storedMetrics[0].Name)
	assert.Equal(t, "counter", storedMetrics[0].Type)
	assert.Equal(t, delta, *storedMetrics[0].Delta)
}

func TestGetMetricValue(t *testing.T) {
	cleanMetricsDatabase(t, testDB)

	logger := zap.NewExample()
	defer logger.Sync()
	storage := New(testDB, logger)

	var val float64 = 5
	metric := dto.Metric{
		Name:  "test_metric",
		Type:  "gauge",
		Value: &val,
	}

	_, err := storage.SetMetric(metric)
	require.NoError(t, err)

	value, err := storage.GetMetricValue("test_metric", "gauge")
	require.NoError(t, err)
	assert.NotNil(t, value, "metric value should be retrieved")
	assert.Equal(t, val, *value)
}

func TestGetMetricsForHTML(t *testing.T) {
	cleanMetricsDatabase(t, testDB)

	logger := zap.NewExample()
	defer logger.Sync()
	storage := New(testDB, logger)

	metrics := []dto.Metric{
		{Name: "metric1", Type: "counter", Value: new(float64)},
		{Name: "metric2", Type: "gauge", Value: new(float64)},
	}
	for _, metric := range metrics {
		_, err := storage.SetMetric(metric)
		require.NoError(t, err)
	}

	storedMetrics, err := storage.GetMetricsForHTML()
	require.NoError(t, err)
	assert.Len(t, storedMetrics, 2, "should retrieve 2 metrics")
	assert.Equal(t, "metric2", storedMetrics[0].Name)
	assert.Equal(t, "metric1", storedMetrics[1].Name)
}

func cleanMetricsDatabase(t *testing.T, db *pgxpool.Pool) {
	_, err := db.Exec(context.Background(), "DELETE FROM metrics")
	require.NoError(t, err)
}
