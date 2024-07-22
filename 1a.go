package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockApplistRepository is a mock of ApplistRepository interface
type MockApplistRepository struct {
	mock.Mock
}

func (m *MockApplistRepository) FindAllComp() (model.Applist, error) {
	args := m.Called()
	return args.Get(0).(model.Applist), args.Error(1)
}

func (m *MockApplistRepository) FindServiceMapData(serviceId string) ([]model.ServiceInfoRes, error) {
	args := m.Called(serviceId)
	return args.Get(0).([]model.ServiceInfoRes), args.Error(1)
}

func (m *MockApplistRepository) FindMetricsMetadata(requestedEntityId string, entityTypes []string) (model.MetricMetadataResponse, error) {
	args := m.Called(requestedEntityId, entityTypes)
	return args.Get(0).(model.MetricMetadataResponse), args.Error(1)
}

func (m *MockApplistRepository) CheckServiceMapByServiceId(serviceId string) (model.ServiceInfo, error) {
	args := m.Called(serviceId)
	return args.Get(0).(model.ServiceInfo), args.Error(1)
}

func (m *MockApplistRepository) FindJourneyDetails(params entity.JourneyParams) (entity.JourneyItsmDetailsRes, error) {
	args := m.Called(params)
	return args.Get(0).(entity.JourneyItsmDetailsRes), args.Error(1)
}

func (m *MockApplistRepository) FindJourneyList() ([]entity.JourneyListRes, error) {
	args := m.Called()
	return args.Get(0).([]entity.JourneyListRes), args.Error(1)
}

// NewApplistApi creates a new instance of DefaultApplistApi
func NewApplistApi(repository model.ApplistRepository) *DefaultApplistApi {
	return &DefaultApplistApi{repo: repository}
}

// Test cases

func TestGetAllAppComp(t *testing.T) {
	mockRepo := new(MockApplistRepository)
	expectedResult := model.Applist{ /* initialize with some data */ }
	mockRepo.On("FindAllComp").Return(expectedResult, nil)

	api := NewApplistApi(mockRepo)
	result, err := api.GetAllAppComp()

	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
	mockRepo.AssertExpectations(t)
}

func TestGetServiceMapData(t *testing.T) {
	mockRepo := new(MockApplistRepository)
	serviceId := "test-service-id"
	expectedResult := []model.ServiceInfoRes{ /* initialize with some data */ }
	mockRepo.On("FindServiceMapData", serviceId).Return(expectedResult, nil)

	api := NewApplistApi(mockRepo)
	result, err := api.GetServiceMapData(serviceId)

	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
	mockRepo.AssertExpectations(t)
}

func TestGetMetricsMetadata(t *testing.T) {
	mockRepo := new(MockApplistRepository)
	requestedEntityId := "test-entity-id"
	entityTypes := []string{"type1", "type2"}
	expectedResult := model.MetricMetadataResponse{ /* initialize with some data */ }
	mockRepo.On("FindMetricsMetadata", requestedEntityId, entityTypes).Return(expectedResult, nil)

	api := NewApplistApi(mockRepo)
	result, err := api.GetMetricsMetadata(requestedEntityId, entityTypes)

	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
	mockRepo.AssertExpectations(t)
}

func TestCheckServiceMap(t *testing.T) {
	mockRepo := new(MockApplistRepository)
	serviceId := "test-service-id"
	expectedResult := model.ServiceInfo{ /* initialize with some data */ }
	mockRepo.On("CheckServiceMapByServiceId", serviceId).Return(expectedResult, nil)

	api := NewApplistApi(mockRepo)
	result, err := api.CheckServiceMap(serviceId)

	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
	mockRepo.AssertExpectations(t)
}

func TestGetJourneyDetails(t *testing.T) {
	mockRepo := new(MockApplistRepository)
	params := entity.JourneyParams{ /* initialize with some data */ }
	expectedResult := entity.JourneyItsmDetailsRes{ /* initialize with some data */ }
	mockRepo.On("FindJourneyDetails", params).Return(expectedResult, nil)

	api := NewApplistApi(mockRepo)
	result, err := api.GetJourneyDetails(params)

	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
	mockRepo.AssertExpectations(t)
}

func TestGetJourneyList(t *testing.T) {
	mockRepo := new(MockApplistRepository)
	expectedResult := []entity.JourneyListRes{ /* initialize with some data */ }
	mockRepo.On("FindJourneyList").Return(expectedResult, nil)

	api := NewApplistApi(mockRepo)
	result, err := api.GetJourneyList()

	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
	mockRepo.AssertExpectations(t)
}