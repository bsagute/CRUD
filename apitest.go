package api

import (
	"errors"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.aexp.com/amex-eng/go-paved-road/pkg/model"
)

// MockAppComponentRepository is a mock implementation of the AppComponentRepository interface
type MockAppComponentRepository struct {
	mock.Mock
}

func (m *MockAppComponentRepository) FindAllAppComp(level *url.Values) (model.AppComponent, error) {
	args := m.Called(level)
	return args.Get(0).(model.AppComponent), args.Error(1)
}

func TestGetAllAppComp(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		mockRepo := new(MockAppComponentRepository)
		api := NewAppComponentApi(mockRepo)

		// Create sample data
		sampleLevel := &url.Values{}
		sampleComponent := model.AppComponent{}
		mockRepo.On("FindAllAppComp", sampleLevel).Return(sampleComponent, nil)

		// Call the method
		result, err := api.GetAllAppComp(sampleLevel)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, sampleComponent, result)

		// Ensure that the expectations were met
		mockRepo.AssertExpectations(t)
	})

	t.Run("failure case", func(t *testing.T) {
		mockRepo := new(MockAppComponentRepository)
		api := NewAppComponentApi(mockRepo)

		// Create sample data
		sampleLevel := &url.Values{}
		expectedError := errors.New("some error")
		mockRepo.On("FindAllAppComp", sampleLevel).Return(model.AppComponent{}, expectedError)

		// Call the method
		result, err := api.GetAllAppComp(sampleLevel)

		// Assertions
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Equal(t, model.AppComponent{}, result)

		// Ensure that the expectations were met
		mockRepo.AssertExpectations(t)
	})
}