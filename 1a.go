package api

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/your_module_path/model"
)

// MockMetricGaugeRepository is a mock of the MetricGaugeRepository interface
type MockMetricGaugeRepository struct {
	mock.Mock
}

func (m *MockMetricGaugeRepository) GetMetricGauge(queriesValue *url.Values) (model.MetricGauge, error) {
	args := m.Called(queriesValue)
	return args.Get(0).(model.MetricGauge), args.Error(1)
}

func (m *MockMetricGaugeRepository) GetMetricGaugeJson(gaugeRequest *model.GraphRequest) (model.MetricGauge, error) {
	args := m.Called(gaugeRequest)
	return args.Get(0).(model.MetricGauge), args.Error(1)
}

// NewMetricGaugeApi creates a new instance of DefaultMetricGaugeApi
func NewMetricGaugeApi(repository model.MetricGaugeRepository) DefaultMetricGaugeApi {
	return DefaultMetricGaugeApi{repo: repository}
}

// Test functions

func TestGetMetricGauge(t *testing.T) {
	mockRepo := new(MockMetricGaugeRepository)
	queriesValue := &url.Values{ /* initialize with some data */ }
	expectedResult := model.MetricGauge{ /* initialize with some data */ }
	mockRepo.On("GetMetricGauge", queriesValue).Return(expectedResult, nil)

	api := NewMetricGaugeApi(mockRepo)
	result, err := api.GetMetricGauge(queriesValue)

	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
	mockRepo.AssertExpectations(t)
}

func TestGetMetricGaugeJson(t *testing.T) {
	mockRepo := new(MockMetricGaugeRepository)
	gaugeRequest := &model.GraphRequest{ /* initialize with some data */ }
	expectedResult := model.MetricGauge{ /* initialize with some data */ }
	mockRepo.On("GetMetricGaugeJson", gaugeRequest).Return(expectedResult, nil)

	api := NewMetricGaugeApi(mockRepo)
	result, err := api.GetMetricGaugeJson(gaugeRequest)

	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
	mockRepo.AssertExpectations(t)
}