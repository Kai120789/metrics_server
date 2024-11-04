package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"server/internal/config"
	"server/internal/dto"
	"server/internal/models"
	"server/internal/transport/http/handler"
	"server/internal/utils"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type MockHandlerer struct {
	mock.Mock
}

func (m *MockHandlerer) SetUpdates(metrics []dto.Metric) ([]models.Metric, error) {
	args := m.Called(metrics)
	return args.Get(0).([]models.Metric), args.Error(1)
}

func (m *MockHandlerer) SetMetric(metric dto.Metric) (*models.Metric, error) {
	args := m.Called(metric)
	return args.Get(0).(*models.Metric), args.Error(1)
}

func (m *MockHandlerer) GetMetricValue(name, typeStr string) (*int64, error) {
	args := m.Called(name, typeStr)
	return args.Get(0).(*int64), args.Error(1)
}

func (m *MockHandlerer) GetHTML(w http.ResponseWriter) error {
	args := m.Called(w)
	return args.Error(0)
}

func TestSetUpdates(t *testing.T) {
	// init mock for handler and injections
	mockService := new(MockHandlerer)
	logger := zap.NewNop()
	conf := &config.Config{SecretKey: "test-secret"}
	h := handler.New(mockService, logger, conf)

	// metrics for send and res
	metrics := []dto.Metric{
		{Name: "metric1", Type: "gauge", Value: new(float64)},
		{Name: "metric2", Type: "counter", Delta: new(int64)},
	}
	reqBody, _ := json.Marshal(metrics)

	// expected result
	expectedMetrics := []models.Metric{}

	mockService.On("SetUpdates", mock.MatchedBy(func(arg []dto.Metric) bool {
		// check is value match
		return len(arg) == len(metrics) &&
			arg[0].Name == metrics[0].Name &&
			arg[0].Type == metrics[0].Type &&
			arg[1].Name == metrics[1].Name &&
			arg[1].Type == metrics[1].Type
	})).Return(expectedMetrics, nil)

	// gen hash
	expectedHash := utils.GenerateHash(conf.SecretKey)
	req := httptest.NewRequest(http.MethodPost, "/updates", bytes.NewBuffer(reqBody))
	req.Header.Set("Hash", expectedHash) // compare hash
	w := httptest.NewRecorder()

	h.SetUpdates(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	// check status code
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestSetMetric(t *testing.T) {
	// init mock for handler and injections
	mockService := new(MockHandlerer)
	logger := zap.NewNop()
	conf := &config.Config{}
	h := handler.New(mockService, logger, conf)

	// init router
	r := chi.NewRouter()
	r.Post("/{type}/{name}/{value}", h.SetMetric)

	value := 10.0
	metric := dto.Metric{
		Name:  "test_metric",
		Type:  "gauge",
		Value: &value,
	}
	mockService.On("SetMetric", metric).Return(&models.Metric{}, nil)

	req := httptest.NewRequest(http.MethodPost, "/gauge/test_metric/10", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	mockService.AssertCalled(t, "SetMetric", metric)
}

func TestGetMetricValue(t *testing.T) {
	mockService := new(MockHandlerer)
	logger := zap.NewNop()
	conf := &config.Config{}
	h := handler.New(mockService, logger, conf)

	r := chi.NewRouter()
	r.Get("/value/{type}/{name}", h.GetMetricValue)

	value := int64(100)
	mockService.On("GetMetricValue", "test_metric", "gauge").Return(&value, nil)

	req := httptest.NewRequest(http.MethodGet, "/value/gauge/test_metric", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	mockService.AssertCalled(t, "GetMetricValue", "test_metric", "gauge")
}

func TestGetHTML(t *testing.T) {
	mockService := new(MockHandlerer)
	logger := zap.NewNop()
	conf := &config.Config{}
	h := handler.New(mockService, logger, conf)

	mockService.On("GetHTML", mock.Anything).Return(nil)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	h.GetHTML(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mockService.AssertCalled(t, "GetHTML", mock.Anything)
}
