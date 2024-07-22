package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/your_module_path/model"
)

// MockMetricBarColRepository is a mock of the MetricBarColRepository interface
type MockMetricBarColRepository struct {
	mock.Mock
}

func (m *MockMetricBarColRepository) GetMetricBarCol() (model.MetricBarCol, error) {
	args := m.Called()
	return args.Get(0).(model.MetricBarCol), args.Error(1)
}

// NewMetricBarColApi creates a new instance of DefaultMetricBarColApi
func NewMetricBarColApi(repository model.MetricBarColRepository) DefaultMetricBarColApi {
	return DefaultMetricBarColApi{repo: repository}
}

// Test functions

func TestGetMetricBarCol(t *testing.T) {
	mockRepo := new(MockMetricBarColRepository)
	expectedResult := model.MetricBarCol{ /* initialize with some data */ }
	mockRepo.On("GetMetricBarCol").Return(expectedResult, nil)

	api := NewMetricBarColApi(mockRepo)
	result, err := api.GetMetricBarCol()

	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
	mockRepo.AssertExpectations(t)
}