package api

import (
    "net/url"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "github.com/yourusername/yourproject/pkg/model" // Adjust the import path accordingly
)

// MockMetricTableRepository is a mock implementation of the MetricTableRepository interface
type MockMetricTableRepository struct {
    mock.Mock
}

func (m *MockMetricTableRepository) GetMetricTable(id *url.Values) (model.MetricTable, error) {
    args := m.Called(id)
    return args.Get(0).(model.MetricTable), args.Error(1)
}

func (m *MockMetricTableRepository) GetMetricTableJson(graphRequest *model.GraphRequest) (model.MetricTable, error) {
    args := m.Called(graphRequest)
    return args.Get(0).(model.MetricTable), args.Error(1)
}

func (m *MockMetricTableRepository) GetMetricTableMetric(graphRequest *model.GraphRequest) (model.MetricTable, error) {
    args := m.Called(graphRequest)
    return args.Get(0).(model.MetricTable), args.Error(1)
}

func TestGetMetricTable_Success(t *testing.T) {
    // Arrange
    id := &url.Values{}
    expectedMetricTable := model.MetricTable{
        // Fill with appropriate fields
    }

    mockRepo := new(MockMetricTableRepository)
    mockRepo.On("GetMetricTable", id).Return(expectedMetricTable, nil)

    api := DefaultMetricTableApi{
        repo: mockRepo,
    }

    // Act
    result, err := api.GetMetricTable(id)

    // Assert
    assert.NoError(t, err)
    assert.Equal(t, expectedMetricTable, result)
    mockRepo.AssertExpectations(t)
}

func TestGetMetricTable_Error(t *testing.T) {
    // Arrange
    id := &url.Values{}
    expectedError := errors.New("some error")

    mockRepo := new(MockMetricTableRepository)
    mockRepo.On("GetMetricTable", id).Return(model.MetricTable{}, expectedError)

    api := DefaultMetricTableApi{
        repo: mockRepo,
    }

    // Act
    result, err := api.GetMetricTable(id)

    // Assert
    assert.Error(t, err)
    assert.Equal(t, expectedError, err)
    assert.Equal(t, model.MetricTable{}, result)
    mockRepo.AssertExpectations(t)
}

func TestGetMetricTableJson_Success(t *testing.T) {
    // Arrange
    graphRequest := &model.GraphRequest{}
    expectedMetricTable := model.MetricTable{
        // Fill with appropriate fields
    }

    mockRepo := new(MockMetricTableRepository)
    mockRepo.On("GetMetricTableJson", graphRequest).Return(expectedMetricTable, nil)

    api := DefaultMetricTableApi{
        repo: mockRepo,
    }

    // Act
    result, err := api.GetMetricTableJson(graphRequest)

    // Assert
    assert.NoError(t, err)
    assert.Equal(t, expectedMetricTable, result)
    mockRepo.AssertExpectations(t)
}

func TestGetMetricTableJson_Error(t *testing.T) {
    // Arrange
    graphRequest := &model.GraphRequest{}
    expectedError := errors.New("some error")

    mockRepo := new(MockMetricTableRepository)
    mockRepo.On("GetMetricTableJson", graphRequest).Return(model.MetricTable{}, expectedError)

    api := DefaultMetricTableApi{
        repo: mockRepo,
    }

    // Act
    result, err := api.GetMetricTableJson(graphRequest)

    // Assert
    assert.Error(t, err)
    assert.Equal(t, expectedError, err)
    assert.Equal(t, model.MetricTable{}, result)
    mockRepo.AssertExpectations(t)
}

func TestGetMetricTableMetric_Success(t *testing.T) {
    // Arrange
    graphRequest := &model.GraphRequest{}
    expectedMetricTable := model.MetricTable{
        // Fill with appropriate fields
    }

    mockRepo := new(MockMetricTableRepository)
    mockRepo.On("GetMetricTableMetric", graphRequest).Return(expectedMetricTable, nil)

    api := DefaultMetricTableApi{
        repo: mockRepo,
    }

    // Act
    result, err := api.GetMetricTableMetric(graphRequest)

    // Assert
    assert.NoError(t, err)
    assert.Equal(t, expectedMetricTable, result)
    mockRepo.AssertExpectations(t)
}

func TestGetMetricTableMetric_Error(t *testing.T) {
    // Arrange
    graphRequest := &model.GraphRequest{}
    expectedError := errors.New("some error")

    mockRepo := new(MockMetricTableRepository)
    mockRepo.On("GetMetricTableMetric", graphRequest).Return(model.MetricTable{}, expectedError)

    api := DefaultMetricTableApi{
        repo: mockRepo,
    }

    // Act
    result, err := api.GetMetricTableMetric(graphRequest)

    // Assert
    assert.Error(t, err)
    assert.Equal(t, expectedError, err)
    assert.Equal(t, model.MetricTable{}, result)
    mockRepo.AssertExpectations(t)
}