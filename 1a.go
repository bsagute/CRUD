// repository_test.go
package main

import (
    "errors"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "gorm.io/gorm"
)

// MockGormDB is a mock of GormDB
type MockGormDB struct {
    mock.Mock
}

// Mock function to replace GetAllEntities
func (m *MockGormDB) GetAllEntities(db *gorm.DB) ([]Entity, error) {
    args := m.Called(db)
    return args.Get(0).([]Entity), args.Error(1)
}

// Inject the mock function
func GetAllEntities(db *gorm.DB) ([]Entity, error) {
    mockDB := &MockGormDB{}
    return mockDB.GetAllEntities(db)
}

func TestAllAppinfoRepositoryDb(t *testing.T) {
    var gormDB *gorm.DB
    mockDB := new(MockGormDB)

    t.Run("should return app info when no error", func(t *testing.T) {
        entities := []Entity{
            {EntityId: 1, EntityName: "App1", EntityDescription: "Description1", MetricInfos: "Metrics1"},
            {EntityId: 2, EntityName: "App2", EntityDescription: "Description2", MetricInfos: "Metrics2"},
        }
        mockDB.On("GetAllEntities", gormDB).Return(entities, nil)

        result := AllAppinfoRepositoryDb(gormDB)

        expectedAppInfo := []Appinfo{
            {App_id: 1, App_name: "App1", App_description: "Description1", Components: "Metrics1"},
            {App_id: 2, App_name: "App2", App_description: "Description2", Components: "Metrics2"},
        }
        assert.Equal(t, expectedAppInfo, result.Appinfo)
        mockDB.AssertExpectations(t)
    })

    t.Run("should log fatal on error", func(t *testing.T) {
        mockDB.On("GetAllEntities", gormDB).Return(nil, errors.New("some error"))

        defer func() {
            if r := recover(); r == nil {
                t.Errorf("Expected log fatal but did not occur")
            }
        }()
        AllAppinfoRepositoryDb(gormDB)
    })
}
