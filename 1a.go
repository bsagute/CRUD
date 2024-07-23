package api

import (
    "errors"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "testing"
    "github.com/yourusername/yourproject/pkg/model" // Adjust the import path accordingly
)

// MockMetricReqRepository is a mock implementation of the MetricReqRepository interface
type MockMetricReqRepository struct {
    mock.Mock
}

func (m *MockMetricReqRepository) GetMetricReq() (model.MetricReq, error) {
    args := m.Called()
    return args.Get(0).(model.MetricReq), args.Error(1)
}

func TestGetMetricReq_Success(t *testing.T) {
    // Arrange
    expectedMetricReq := model.MetricReq{
        // Fill with appropriate fields
    }

    mockRepo := new(MockMetricReqRepository)
    mockRepo.On("GetMetricReq").Return(expectedMetricReq, nil)

    api := DefaultMetricReqApi{
        repo: mockRepo,
    }

    // Act
    result, err := api.GetMetricReq()

    // Assert
    assert.NoError(t, err)
    assert.Equal(t, expectedMetricReq, result)
    mockRepo.AssertExpectations(t)
}

func TestGetMetricReq_Error(t *testing.T) {
    // Arrange
    expectedError := errors.New("some error")

    mockRepo := new(MockMetricReqRepository)
    mockRepo.On("GetMetricReq").Return(model.MetricReq{}, expectedError)

    api := DefaultMetricReqApi{
        repo: mockRepo,
    }

    // Act
    result, err := api.GetMetricReq()

    // Assert
    assert.Error(t, err)
    assert.Equal(t, expectedError, err)
    assert.Equal(t, model.MetricReq{}, result)
    mockRepo.AssertExpectations(t)
}