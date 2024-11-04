package memstorage_test

import (
	"testing"

	"server/internal/dto"
	"server/internal/models"
	"server/internal/storage/memstorage"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// Mock for the Logger
type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Info(msg string, fields ...zap.Field) {
	m.Called(msg, fields)
}

func (m *MockLogger) Error(msg string, fields ...zap.Field) {
	m.Called(msg, fields)
}

func TestSetUpdates(t *testing.T) {
	// Arrange
	logger := zap.NewNop() // Use a no-op logger for simplicity
	storage := memstorage.New([]models.Metric{}, logger)

	metrics := []dto.Metric{
		{Name: "test_metric1", Type: "counter", Value: nil, Delta: (new(int64))},
		{Name: "test_metric2", Type: "gauge", Value: nil, Delta: nil},
	}

	// Act
	updatedMetrics, err := storage.SetUpdates(metrics)

	// Assert
	require.NoError(t, err)
	assert.Len(t, updatedMetrics, 2)
	assert.Equal(t, "test_metric1", updatedMetrics[0].Name)
	assert.Equal(t, "test_metric2", updatedMetrics[1].Name)
	assert.Equal(t, uint(1), updatedMetrics[0].ID)
	assert.Equal(t, uint(2), updatedMetrics[1].ID)

	// Check that Delta was updated
	assert.Equal(t, int64(0), *updatedMetrics[0].Delta)
}

func TestSetMetric(t *testing.T) {
	// Arrange
	logger := zap.NewNop()
	storage := memstorage.New([]models.Metric{}, logger)

	value := float64(10)
	metric := dto.Metric{Name: "test_metric", Type: "gauge", Value: &value, Delta: nil}

	// Act
	returnedMetric, err := storage.SetMetric(metric)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "test_metric", returnedMetric.Name)
	assert.Equal(t, "gauge", returnedMetric.Type)
	assert.Equal(t, value, *returnedMetric.Value)
}

func TestGetMetricValue(t *testing.T) {
	// Arrange
	logger := zap.NewNop()
	value := float64(20)
	metric := models.Metric{Name: "test_metric", Type: "gauge", Value: &value, Delta: nil}
	storage := memstorage.New([]models.Metric{metric}, logger)

	// Act
	result, err := storage.GetMetricValue("test_metric", "gauge")

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(value), *result)
}

func TestGetMetricsForHTML(t *testing.T) {
	// Arrange
	logger := zap.NewNop()
	value1 := float64(15)
	value2 := float64(25)
	metrics := []models.Metric{
		{Name: "metric1", Type: "gauge", Value: &value1, Delta: nil},
		{Name: "metric2", Type: "gauge", Value: &value2, Delta: nil},
	}
	storage := memstorage.New(metrics, logger)

	// Act
	returnedMetrics, err := storage.GetMetricsForHTML()

	// Assert
	require.NoError(t, err)
	assert.Len(t, returnedMetrics, 2)
	assert.Equal(t, "metric1", returnedMetrics[0].Name)
	assert.Equal(t, "metric2", returnedMetrics[1].Name)
}
