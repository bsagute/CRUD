// app_info_db_test.go
package db

import (
    "errors"
    "fmt"
    "testing"

    "github.com/google/uuid"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "gorm.io/gorm"
    "github.aexp.com/amex-eng/go-paved-road/pkg/entity"
    "github.aexp.com/amex-eng/go-paved-road/pkg/model"
)

// MockGormDB is a mock of GormDB
type MockGormDB struct {
    mock.Mock
}

// Mock function to replace GetAllEntities
func (m *MockGormDB) GetAllEntities(db *gorm.DB) ([]entity.EntityInfo, error) {
    args := m.Called(db)
    return args.Get(0).([]entity.EntityInfo), args.Error(1)
}

// Override GetAllEntities to use the mock
func GetAllEntities(db *gorm.DB) ([]entity.EntityInfo, error) {
    mockDB := &MockGormDB{}
    return mockDB.GetAllEntities(db)
}

func TestAllAppinfoRepositoryDb(t *testing.T) {
    var gormDB *gorm.DB
    mockDB := new(MockGormDB)

    t.Run("should return app info when no error", func(t *testing.T) {
        entities := []entity.EntityInfo{
            {EntityId: uuid.New(), EntityName: "App1", EntityDescription: "Description1"},
            {EntityId: uuid.New(), EntityName: "App2", EntityDescription: "Description2"},
        }
        mockDB.On("GetAllEntities", gormDB).Return(entities, nil)

        result := AllAppinfoRepositoryDb(gormDB)

        expectedAppInfo := []model.Appinfo{
            {App_name: "App1", App_id: entities[0].EntityId},
            {App_name: "App2", App_id: entities[1].EntityId},
        }
        assert.Equal(t, expectedAppInfo, result.Appinfo)
        mockDB.AssertExpectations(t)
    })

    t.Run("should log fatal on error", func(t *testing.T) {
        mockDB.On("GetAllEntities", gormDB).Return(nil, errors.New("some error"))

        // Mock the log.Fatalln to prevent actual fatal error
        origLogFatal := logFatalln
        defer func() { logFatalln = origLogFatal }()
        logFatalCalled := false
        logFatalln = func(v ...interface{}) {
            logFatalCalled = true
            fmt.Println(v...)
        }

        AllAppinfoRepositoryDb(gormDB)
        assert.True(t, logFatalCalled, "Expected log fatal but did not occur")
    })
}

// Override log.Fatalln to prevent actual fatal error during testing
var logFatalln = func(v ...interface{}) {
    fmt.Println(v...)
}
