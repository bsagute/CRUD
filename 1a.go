package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/your_module_path/model"
)

// MockMetricGraphMultilineRepository is a mock of the MetricGraphMultilineRepository interface
type MockMetricGraphMultilineRepository struct {
	mock.Mock
}

func (m *MockMetricGraphMultilineRepository) GetMetricGraphMultiline() (model.MetricGraphMultiline, error) {
	args := m.Called()
	return args.Get(0).(model.MetricGraphMultiline), args.Error(1)
}

func (m *MockMetricGraphMultilineRepository) GetMetricGraphMultilineJson(graphRequest *model.GraphRequest) (model.MetricGraphMultiline, error) {
	args := m.Called(graphRequest)
	return args.Get(0).(model.MetricGraphMultiline), args.Error(1)
}

func (m *MockMetricGraphMultilineRepository) GetMetricGraphMultilineMetric(graphRequest *model.GraphRequest) (model.MetricGraphMultiline, error) {
	args := m.Called(graphRequest)
	return args.Get(0).(model.MetricGraphMultiline), args.Error(1)
}

func (m *MockMetricGraphMultilineRepository) GetMetricGraphMultilineDynamic(graphRequest *model.GraphRequest) (model.PodMetricGraphMultiline, error) {
	args := m.Called(graphRequest)
	return args.Get(0).(model.PodMetricGraphMultiline), args.Error(1)
}

// NewMetricGraphMultilineApi creates a new instance of DefaultMetricGraphMultilineApi
func NewMetricGraphMultilineApi(repository model.MetricGraphMultilineRepository) DefaultMetricGraphMultilineApi {
	return DefaultMetricGraphMultilineApi{repo: repository}
}

// Test functions

func TestGetMetricGraphMultiline(t *testing.T) {
	mockRepo := new(MockMetricGraphMultilineRepository)
	expectedResult := model.MetricGraphMultiline{ /* initialize with some data */ }
	mockRepo.On("GetMetricGraphMultiline").Return(expectedResult, nil)

	api := NewMetricGraphMultilineApi(mockRepo)
	result, err := api.GetMetricGraphMultiline()

	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
	mockRepo.AssertExpectations(t)
}

func TestGetMetricGraphMultilineJson(t *testing.T) {
	mockRepo := new(MockMetricGraphMultilineRepository)
	graphRequest := &model.GraphRequest{ /* initialize with some data */ }
	expectedResult := model.MetricGraphMultiline{ /* initialize with some data */ }
	mockRepo.On("GetMetricGraphMultilineJson", graphRequest).Return(expectedResult, nil)

	api := NewMetricGraphMultilineApi(mockRepo)
	result, err := api.GetMetricGraphMultilineJson(graphRequest)

	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
	mockRepo.AssertExpectations(t)
}

func TestGetMetricGraphMultilineMetric(t *testing.T) {
	mockRepo := new(MockMetricGraphMultilineRepository)
	graphRequest := &model.GraphRequest{ /* initialize with some data */ }
	expectedResult := model.MetricGraphMultiline{ /* initialize with some data */ }
	mockRepo.On("GetMetricGraphMultilineMetric", graphRequest).Return(expectedResult, nil)

	api := NewMetricGraphMultilineApi(mockRepo)
	result, err := api.GetMetricGraphMultilineMetric(graphRequest)

	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
	mockRepo.AssertExpectations(t)
}

func TestGetMetricGraphMultilineDynamic(t *testing.T) {
	mockRepo := new(MockMetricGraphMultilineRepository)
	graphRequest := &model.GraphRequest{ /* initialize with some data */ }
	expectedResult := model.PodMetricGraphMultiline{ /* initialize with some data */ }
	mockRepo.On("GetMetricGraphMultilineDynamic", graphRequest).Return(expectedResult, nil)

	api := NewMetricGraphMultilineApi(mockRepo)
	result, err := api.GetMetricGraphMultilineDynamic(graphRequest)

	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
	mockRepo.AssertExpectations(t)
}