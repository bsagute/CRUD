package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock repository to be used in tests
type MockMetricGraphMultilineRepository struct {
	mock.Mock
}

func (m *MockMetricGraphMultilineRepository) GetMetricGraphMultiline() (MetricGraphMultiline, error) {
	args := m.Called()
	return args.Get(0).(MetricGraphMultiline), args.Error(1)
}

func (m *MockMetricGraphMultilineRepository) GetMetricGraphMultilineJson(graphRequest *GraphRequest) (MetricGraphMultiline, error) {
	args := m.Called(graphRequest)
	return args.Get(0).(MetricGraphMultiline), args.Error(1)
}

func (m *MockMetricGraphMultilineRepository) GetMetricGraphMultilineMetric(graphRequest *GraphRequest) (MetricGraphMultiline, error) {
	args := m.Called(graphRequest)
	return args.Get(0).(MetricGraphMultiline), args.Error(1)
}

func (m *MockMetricGraphMultilineRepository) GetMetricGraphMultilineDynamic(graphRequest *GraphRequest) (PodMetricGraphMultiline, error) {
	args := m.Called(graphRequest)
	return args.Get(0).(PodMetricGraphMultiline), args.Error(1)
}

func TestDefaultMetricGraphMultilineApi_GetMetricGraphMultiline(t *testing.T) {
	mockRepo := new(MockMetricGraphMultilineRepository)
	expectedResult := MetricGraphMultiline{}
	mockRepo.On("GetMetricGraphMultiline").Return(expectedResult, nil)

	api := NewMetricGraphMultilineApi(mockRepo)

	result, err := api.GetMetricGraphMultiline()

	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
	mockRepo.AssertExpectations(t)
}

func TestDefaultMetricGraphMultilineApi_GetMetricGraphMultilineJson(t *testing.T) {
	mockRepo := new(MockMetricGraphMultilineRepository)
	graphRequest := &GraphRequest{}
	expectedResult := MetricGraphMultiline{}
	mockRepo.On("GetMetricGraphMultilineJson", graphRequest).Return(expectedResult, nil)

	api := NewMetricGraphMultilineApi(mockRepo)

	result, err := api.GetMetricGraphMultilineJson(graphRequest)

	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
	mockRepo.AssertExpectations(t)
}

func TestDefaultMetricGraphMultilineApi_GetMetricGraphMultilineMetric(t *testing.T) {
	mockRepo := new(MockMetricGraphMultilineRepository)
	graphRequest := &GraphRequest{}
	expectedResult := MetricGraphMultiline{}
	mockRepo.On("GetMetricGraphMultilineMetric", graphRequest).Return(expectedResult, nil)

	api := NewMetricGraphMultilineApi(mockRepo)

	result, err := api.GetMetricGraphMultilineMetric(graphRequest)

	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
	mockRepo.AssertExpectations(t)
}

func TestDefaultMetricGraphMultilineApi_GetMetricGraphMultilineDynamic(t *testing.T) {
	mockRepo := new(MockMetricGraphMultilineRepository)
	graphRequest := &GraphRequest{}
	expectedResult := PodMetricGraphMultiline{}
	mockRepo.On("GetMetricGraphMultilineDynamic", graphRequest).Return(expectedResult, nil)

	api := NewMetricGraphMultilineApi(mockRepo)

	result, err := api.GetMetricGraphMultilineDynamic(graphRequest)

	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
	mockRepo.AssertExpectations(t)
}