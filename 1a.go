package api

import (
	"errors"
	"net/url"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/your/repository/model"
)

// Mock for the AppinfoRepository
type MockAppinfoRepository struct {
	mock.Mock
}

func (m *MockAppinfoRepository) FindAll(params *url.Values) ([]model.Appinfo, error) {
	args := m.Called(params)
	return args.Get(0).([]model.Appinfo), args.Error(1)
}

func TestGetAllAppInfo(t *testing.T) {
	mockRepo := new(MockAppinfoRepository)
	api := NewAppinfoApi(mockRepo)

	params := &url.Values{}
	expectedResult := []model.Appinfo{
		{
			App_id:         uuid.New(),
			App_name:       "App1",
			App_type:       "Type1",
			App_description: "Description1",
			Components:     nil,
			Metricinfo:     nil,
		},
		{
			App_id:         uuid.New(),
			App_name:       "App2",
			App_type:       "Type2",
			App_description: "Description2",
			Components:     nil,
			Metricinfo:     nil,
		},
	}

	// Setting up expectations
	mockRepo.On("FindAll", params).Return(expectedResult, nil)

	// Call the method
	result, err := api.GetAllAppInfo(params)

	// Assert the expectations
	assert.Nil(t, err)
	assert.Equal(t, expectedResult, result)

	mockRepo.AssertExpectations(t)
}

func TestGetAllAppInfo_Error(t *testing.T) {
	mockRepo := new(MockAppinfoRepository)
	api := NewAppinfoApi(mockRepo)

	params := &url.Values{}
	expectedError := errors.New("some error")

	// Setting up expectations
	mockRepo.On("FindAll", params).Return(nil, expectedError)

	// Call the method
	result, err := api.GetAllAppInfo(params)

	// Assert the expectations
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)

	mockRepo.AssertExpectations(t)
}