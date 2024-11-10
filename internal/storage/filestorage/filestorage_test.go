package filestorage_test

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"server/internal/dto"
	"server/internal/models"
	"server/internal/storage/filestorage"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// Создайте мок для FileSystem
type MockFileSystem struct {
	mock.Mock
}

func TestMain(m *testing.M) {
	if _, err := os.Stat("./metrics.json"); os.IsNotExist(err) {
		file, err := os.Create("./metrics.json")
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		defer file.Close()
		fmt.Println("file succesfully created:", "./metrics.json")

	} else {
		fmt.Println("file is exist now:", "./metrics.json")
	}

	m.Run()
}

func (m *MockFileSystem) ReadFile(name string) ([]byte, error) {
	args := m.Called(name)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockFileSystem) WriteFile(name string, data []byte, perm os.FileMode) error {
	args := m.Called(name, data, perm)
	return args.Error(0)
}

func (m *MockFileSystem) Create(name string) (*os.File, error) {
	args := m.Called(name)
	return args.Get(0).(*os.File), args.Error(1)
}

func (m *MockFileSystem) Stat(name string) (os.FileInfo, error) {
	args := m.Called(name)
	return args.Get(0).(os.FileInfo), args.Error(1)
}

// Примеры тестов
func TestSetUpdates(t *testing.T) {
	mockFileSystem := new(MockFileSystem)
	logger := zap.NewNop()
	storage := filestorage.New("./metrics.json", logger)

	var delta int64 = 5

	// Arrange
	metrics := []dto.Metric{
		{Name: "test_metric1", Type: "counter", Value: nil, Delta: &delta},
		{Name: "test_metric2", Type: "gauge", Value: new(float64), Delta: nil},
	}

	// Act
	returnedMetrics, err := storage.SetUpdates(metrics)

	// Assert
	require.NoError(t, err)
	assert.Len(t, returnedMetrics, 2)
	assert.Equal(t, "test_metric1", returnedMetrics[0].Name)
	assert.Equal(t, "test_metric2", returnedMetrics[1].Name)

	// Verify that the mock expectations were met
	mockFileSystem.AssertExpectations(t)
}

func TestSetMetric(t *testing.T) {
	mockFileSystem := new(MockFileSystem)
	logger := zap.NewNop()
	storage := filestorage.New("./metrics.json", logger)

	// Arrange
	value := 10.0
	metric := dto.Metric{Name: "test_metric", Type: "gauge", Value: &value}

	// Act
	returnedMetric, err := storage.SetMetric(metric)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "test_metric", returnedMetric.Name)
	assert.Equal(t, "gauge", returnedMetric.Type)

	// Verify that the mock expectations were met
	mockFileSystem.AssertExpectations(t)
}

func TestGetMetricValue(t *testing.T) {
	mockFileSystem := new(MockFileSystem)
	logger := zap.NewNop()
	storage := filestorage.New("./metrics.json", logger)

	// Arrange
	metric := []models.Metric{
		{Name: "test_metric", Type: "gauge", Value: nil},
	}

	data, _ := json.Marshal(metric)

	_ = data

	// Act
	value, err := storage.GetMetricValue("test_metric", "gauge")

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, value)

	// Verify that the mock expectations were met
	mockFileSystem.AssertExpectations(t)
}
