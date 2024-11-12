package service_test

import (
	"net/http"
	"net/http/httptest"
	"server/internal/dto"
	"server/internal/models"
	"server/internal/service"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type MockStorager struct {
	mock.Mock
}

func (m *MockStorager) SetUpdates(metrics []dto.Metric) ([]models.Metric, error) {
	args := m.Called(metrics)
	return args.Get(0).([]models.Metric), args.Error(1)
}

func (m *MockStorager) SetMetric(metric dto.Metric) (*models.Metric, error) {
	args := m.Called(metric)
	return args.Get(0).(*models.Metric), args.Error(1)
}

func (m *MockStorager) GetMetricValue(name string, typeStr string) (*float64, error) {
	args := m.Called(name, typeStr)
	return args.Get(0).(*float64), args.Error(1)
}

func (m *MockStorager) GetMetricsForHTML() ([]models.Metric, error) {
	args := m.Called()
	return args.Get(0).([]models.Metric), args.Error(1)
}

func TestSetUpdates(t *testing.T) {
	mockStorage := new(MockStorager)
	srv := service.New(mockStorage, &zap.Logger{})

	metrics := []dto.Metric{
		{Name: "test_metric1", Type: "counter", Value: nil, Delta: new(int64)},
	}

	expectedMetrics := []models.Metric{
		{ID: 1, Name: "test_metric1", Type: "counter", Value: nil, Delta: new(int64), CreatedAt: time.Now()},
	}

	mockStorage.On("SetUpdates", metrics).Return(expectedMetrics, nil)

	result, err := srv.SetUpdates(metrics)

	assert.NoError(t, err)
	assert.Equal(t, expectedMetrics, result)
	mockStorage.AssertExpectations(t)
}

func TestSetMetric(t *testing.T) {
	mockStorage := new(MockStorager)
	srv := service.New(mockStorage, &zap.Logger{})

	metric := dto.Metric{Name: "test_metric", Type: "gauge", Value: new(float64), Delta: nil}
	expectedMetric := &models.Metric{ID: 1, Name: "test_metric", Type: "gauge", Value: new(float64), CreatedAt: time.Now()}

	mockStorage.On("SetMetric", metric).Return(expectedMetric, nil)

	result, err := srv.SetMetric(metric)

	assert.NoError(t, err)
	assert.Equal(t, expectedMetric, result)
	mockStorage.AssertExpectations(t)
}

func TestGetMetricValue(t *testing.T) {
	mockStorage := new(MockStorager)
	srv := service.New(mockStorage, &zap.Logger{})

	expectedValue := float64(42.44)
	mockStorage.On("GetMetricValue", "test_metric", "counter").Return(&expectedValue, nil)

	value, err := srv.GetMetricValue("test_metric", "counter")

	assert.NoError(t, err)
	assert.Equal(t, &expectedValue, value)
	mockStorage.AssertExpectations(t)
}

func TestGetHTML(t *testing.T) {
	mockStorage := new(MockStorager)
	srv := service.New(mockStorage, &zap.Logger{})

	metrics := []models.Metric{
		{ID: 1, Name: "test_metric1", Type: "counter", Value: new(float64), Delta: new(int64), CreatedAt: time.Now()},
	}

	mockStorage.On("GetMetricsForHTML").Return(metrics, nil)

	recorder := httptest.NewRecorder()

	err := srv.GetHTML(recorder)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "<h1>Metrics</h1>")
	assert.Contains(t, recorder.Body.String(), "<td>test_metric1</td>")

	mockStorage.AssertExpectations(t)
}
