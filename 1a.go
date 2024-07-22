package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/your_module_path/model"
)

// MockItsmMetricRepository is a mock of the ItsmMetricRepository interface
type MockItsmMetricRepository struct {
	mock.Mock
}

func (m *MockItsmMetricRepository) GetItsmMetric() (model.ItsmMetric, error) {
	args := m.Called()
	return args.Get(0).(model.ItsmMetric), args.Error(1)
}

func (m *MockItsmMetricRepository) GetItsmMetricJson(graphRequest *model.GraphRequest) (model.ItsmMetric, error) {
	args := m.Called(graphRequest)
	return args.Get(0).(model.ItsmMetric), args.Error(1)
}

// NewItsmMetricApi creates a new instance of DefaultItsmMetricApi
func NewItsmMetricApi(repository model.ItsmMetricRepository) DefaultItsmMetricApi {
	return DefaultItsmMetricApi{repo: repository}
}

// Test functions

func TestGetItsmMetric(t *testing.T) {
	mockRepo := new(MockItsmMetricRepository)
	expectedResult := model.ItsmMetric{ /* initialize with some data */ }
	mockRepo.On("GetItsmMetric").Return(expectedResult, nil)

	api := NewItsmMetricApi(mockRepo)
	result, err := api.GetItsmMetric()

	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
	mockRepo.AssertExpectations(t)
}

func TestGetItsmMetricJson(t *testing.T) {
	mockRepo := new(MockItsmMetricRepository)
	graphRequest := &model.GraphRequest{ /* initialize with some data */ }
	expectedResult := model.ItsmMetric{ /* initialize with some data */ }
	mockRepo.On("GetItsmMetricJson", graphRequest).Return(expectedResult, nil)

	api := NewItsmMetricApi(mockRepo)
	result, err := api.GetItsmMetricJson(graphRequest)

	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
	mockRepo.AssertExpectations(t)
}