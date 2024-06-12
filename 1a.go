package db_test

import (
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "github.com/stretchr/testify/mock"
    
    "github.com/path/to/your/project/db" // Update this path to the actual import path of your project
    "github.com/path/to/your/project/pkg/model"
)

// Mock the dependencies if needed (like gorm.DB). If not needed, you can remove this part.
type MockDB struct {
    mock.Mock
}

func (m *MockDB) GetItsmMetric() (model.ItsmMetric, error) {
    args := m.Called()
    return args.Get(0).(model.ItsmMetric), args.Error(1)
}

func TestGetItsmMetric(t *testing.T) {
    // Create a new instance of the ItsmMetricRepositoryDb with any mocked dependencies
    repository := db.ItsmMetricRepositoryDb{
        getGorm: &MockDB{},
    }

    expectedItsmMetric := model.ItsmMetric{
        Type: "itsm-metric",
    }

    it, err := repository.GetItsmMetric()

    require.NoError(t, err)
    assert.Equal(t, expectedItsmMetric, it)
}