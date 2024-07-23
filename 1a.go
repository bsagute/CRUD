package db_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"

	"github.com/YOUR_USERNAME/YOUR_PROJECT/db"
	"github.com/YOUR_USERNAME/YOUR_PROJECT/entity"
)

// MockGormDB is a mock type for the Gorm DB
type MockGormDB struct {
	mock.Mock
}

func (m *MockGormDB) Raw(sql string, values ...interface{}) *gorm.DB {
	args := m.Called(sql, values)
	return args.Get(0).(*gorm.DB)
}

func (m *MockGormDB) Scan(dest interface{}) *gorm.DB {
	args := m.Called(dest)
	return args.Get(0).(*gorm.DB)
}

func TestFindJourneyList(t *testing.T) {
	mockDB := new(MockGormDB)
	repo := db.ApplistRepositoryDb{
		getGorm: mockDB,
	}

	expectedResult := []entity.JourneyListRes{
		{JourneyID: "1", JourneyName: "Test Journey"},
	}

	mockDB.On("Raw", "select distinct journey_id, journey_name from journey_carid_mapping").Return(mockDB)
	mockDB.On("Scan", &[]entity.JourneyListRes{}).Run(func(args mock.Arguments) {
		dest := args.Get(0).(*[]entity.JourneyListRes)
		*dest = expectedResult
	}).Return(&gorm.DB{})

	result, err := repo.FindJourneyList()
	assert.Nil(t, err)
	assert.Equal(t, expectedResult, result)
	mockDB.AssertExpectations(t)
}

func TestFindJourneyListError(t *testing.T) {
	mockDB := new(MockGormDB)
	repo := db.ApplistRepositoryDb{
		getGorm: mockDB,
	}

	mockDB.On("Raw", "select distinct journey_id, journey_name from journey_carid_mapping").Return(mockDB)
	mockDB.On("Scan", &[]entity.JourneyListRes{}).Return(&gorm.DB{
		Error: errors.New("database error"),
	})

	result, err := repo.FindJourneyList()
	assert.NotNil(t, err)
	assert.Nil(t, result)
	mockDB.AssertExpectations(t)
}