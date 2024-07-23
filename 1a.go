package api

import (
    "errors"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

// Mock repository
type MockMetricGraphMultiLineRepository struct {
    mock.Mock
}

func (m *MockMetricGraphMultiLineRepository) GetMetricGraphMultiLine() (model.MetricGraphMultiLine, error) {
    args := m.Called()
    return args.Get(0).(model.MetricGraphMultiLine), args.Error(1)
}

func (m *MockMetricGraphMultiLineRepository) GetMetricGraphMultiLineJson(graphRequest *model.GraphRequest) (model.MetricGraphMultiLine, error) {
    args := m.Called(graphRequest)
    return args.Get(0).(model.MetricGraphMultiLine), args.Error(1)
}

func (m *MockMetricGraphMultiLineRepository) GetMetricGraphMultiLineMetric(graphRequest *model.GraphRequest) (model.MetricGraphMultiLine, error) {
    args := m.Called(graphRequest)
    return args.Get(0).(model.MetricGraphMultiLine), args.Error(1)
}

func (m *MockMetricGraphMultiLineRepository) GetMetricGraphMultiLineDynamic(graphRequest *model.GraphRequest) (model.PodMetricGraphMultiLine, error) {
    args := m.Called(graphRequest)
    return args.Get(0).(model.PodMetricGraphMultiLine), args.Error(1)
}

func TestDefaultMetricGraphMultiLineApi_GetMetricGraphMultiLine(t *testing.T) {
    mockRepo := new(MockMetricGraphMultiLineRepository)
    service := DefaultMetricGraphMultiLineApi{repo: mockRepo}

    expectedResult := model.MetricGraphMultiLine{}
    mockRepo.On("GetMetricGraphMultiLine").Return(expectedResult, nil)

    result, err := service.GetMetricGraphMultiLine()
    assert.NoError(t, err)
    assert.Equal(t, expectedResult, result)
}

func TestDefaultMetricGraphMultiLineApi_GetMetricGraphMultiLineJson(t *testing.T) {
    mockRepo := new(MockMetricGraphMultiLineRepository)
    service := DefaultMetricGraphMultiLineApi{repo: mockRepo}

    expectedResult := model.MetricGraphMultiLine{}
    graphRequest := &model.GraphRequest{}
    mockRepo.On("GetMetricGraphMultiLineJson", graphRequest).Return(expectedResult, nil)

    result, err := service.GetMetricGraphMultiLineJson(graphRequest)
    assert.NoError(t, err)
    assert.Equal(t, expectedResult, result)
}

func TestDefaultMetricGraphMultiLineApi_GetMetricGraphMultiLineMetric(t *testing.T) {
    mockRepo := new(MockMetricGraphMultiLineRepository)
    service := DefaultMetricGraphMultiLineApi{repo: mockRepo}

    expectedResult := model.MetricGraphMultiLine{}
    graphRequest := &model.GraphRequest{}
    mockRepo.On("GetMetricGraphMultiLineMetric", graphRequest).Return(expectedResult, nil)

    result, err := service.GetMetricGraphMultiLineMetric(graphRequest)
    assert.NoError(t, err)
    assert.Equal(t, expectedResult, result)
}

func TestDefaultMetricGraphMultiLineApi_GetMetricGraphMultiLineDynamic(t *testing.T) {
    mockRepo := new(MockMetricGraphMultiLineRepository)
    service := DefaultMetricGraphMultiLineApi{repo: mockRepo}

    expectedResult := model.PodMetricGraphMultiLine{}
    graphRequest := &model.GraphRequest{}
    mockRepo.On("GetMetricGraphMultiLineDynamic", graphRequest).Return(expectedResult, nil)

    result, err := service.GetMetricGraphMultiLineDynamic(graphRequest)
    assert.NoError(t, err)
    assert.Equal(t, expectedResult, result)
}