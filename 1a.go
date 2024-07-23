package db_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"

	"github.com/YOUR_USERNAME/YOUR_PROJECT/db"
	"github.com/YOUR_USERNAME/YOUR_PROJECT/model"
)

// MockGormDB is a mock type for the Gorm DB
type MockGormDB struct {
	mock.Mock
}

func (m *MockGormDB) Raw(sql string, values ...interface{}) *gorm.DB {
	args := m.Called(sql, values)
	return args.Get(0).(*gorm.DB)
}

func (m *MockGormDB) Model(value interface{}) *gorm.DB {
	args := m.Called(value)
	return args.Get(0).(*gorm.DB)
}

func (m *MockGormDB) Where(query interface{}, args ...interface{}) *gorm.DB {
	mockArgs := m.Called(query, args)
	return mockArgs.Get(0).(*gorm.DB)
}

func (m *MockGormDB) Find(dest interface{}, conds ...interface{}) *gorm.DB {
	args := m.Called(dest, conds)
	return args.Get(0).(*gorm.DB)
}

func TestFindAllComp(t *testing.T) {
	mockDB := new(MockGormDB)
	repo := db.ApplistRepositoryDb{
		getGorm: mockDB,
	}

	expectedResult := model.Applist{}
	mockDB.On("Raw", mock.Anything, mock.Anything).Return(&gorm.DB{})

	result, err := repo.FindAllComp()
	assert.Nil(t, err)
	assert.Equal(t, expectedResult, result)
	mockDB.AssertExpectations(t)
}

func TestFindServiceMapData(t *testing.T) {
	mockDB := new(MockGormDB)
	repo := db.ApplistRepositoryDb{
		getGorm: mockDB,
	}

	expectedResult := []model.ServiceInfoRes{}
	mockDB.On("Raw", mock.Anything, mock.Anything).Return(&gorm.DB{})

	result, err := repo.FindServiceMapData("some-service-id")
	assert.Nil(t, err)
	assert.Equal(t, expectedResult, result)
	mockDB.AssertExpectations(t)
}

func TestFindMetricsMetadata(t *testing.T) {
	repo := db.ApplistRepositoryDb{}

	expectedResult := []model.MetricMetadataResponse{}

	result, err := repo.FindMetricsMetadata("some-entity-id", []string{"type1", "type2"})
	assert.Nil(t, err)
	assert.Equal(t, expectedResult, result)
}

func TestCheckServiceMapByServiceId(t *testing.T) {
	mockDB := new(MockGormDB)
	repo := db.ApplistRepositoryDb{
		getGorm: mockDB,
	}

	expectedResult := model.ServiceInfo{}
	mockDB.On("Model", mock.Anything).Return(mockDB)
	mockDB.On("Where", "service_id = ?", mock.Anything).Return(mockDB)
	mockDB.On("Find", mock.Anything).Return(&gorm.DB{})

	result, err := repo.CheckServiceMapByServiceId("some-service-id")
	assert.Nil(t, err)
	assert.Equal(t, &expectedResult, result)
	mockDB.AssertExpectations(t)
}

func TestFindJourneyDetails(t *testing.T) {
	mockDB := new(MockGormDB)
	repo := db.ApplistRepositoryDb{
		getGorm: mockDB,
	}

	expectedResult := []model.JourneyItsmDetailsRes{}
	mockDB.On("Raw", mock.Anything, mock.Anything).Return(&gorm.DB{})

	params := model.JourneyParams{
		StartDate: "2023-01-01",
		EndDate:   "2023-12-31",
		Limit:     10,
		Offset:    0,
	}
	result, err := repo.FindJourneyDetails(params)
	assert.Nil(t, err)
	assert.Equal(t, expectedResult, result)
	mockDB.AssertExpectations(t)
}

func TestFindJourneyList(t *testing.T) {
	mockDB := new(MockGormDB)
	repo := db.ApplistRepositoryDb{
		getGorm: mockDB,
	}

	expectedResult := []model.JourneyListRes{}
	mockDB.On("Raw", "select distinct journey_id, journey_name from journey_carid_mapping").Return(&gorm.DB{})

	result, err := repo.FindJourneyList()
	assert.Nil(t, err)
	assert.Equal(t, expectedResult, result)
	mockDB.AssertExpectations(t)
}

func TestGetApplistData(t *testing.T) {
	mockDB := new(MockGormDB)
	mockDB.On("Preload", "LayoutsInfos", "graph_type = ?", "table").Return(mockDB)
	mockDB.On("Where", "parent_entity_id is null").Return(mockDB)
	mockDB.On("Find", mock.Anything).Return(&gorm.DB{})

	result, err := db.getApplistData(mockDB)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	mockDB.AssertExpectations(t)
}