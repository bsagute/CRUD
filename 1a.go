package db

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestGetMetricTable(t *testing.T) {
	// Mocking the necessary data and functions
	gormDB := &gorm.DB{} // Mock or initialize a Gorm DB connection
	repo := MetricTableRepositoryDb{getGorm: gormDB}

	query := url.Values{}
	query.Set("someKey", "someValue") // Adjust according to actual query parameters

	expectedMetricTable := model.MetricTable{} // Populate with expected data

	// Mock getTableData to return expected results
	getTableData = func(db *gorm.DB, query url.Values) (model.MetricTable, error) {
		return expectedMetricTable, nil
	}

	metricTable, err := repo.GetMetricTable(query)
	assert.NoError(t, err)
	assert.Equal(t, expectedMetricTable, metricTable)
}

func TestNewMetricTableRepositoryDb(t *testing.T) {
	// Mocking the necessary data and functions
	gormDB := &gorm.DB{} // Mock or initialize a Gorm DB connection

	// Create a mock JSON file
	mockJSON := `{"id": "testID", "data": "testData"}`
	ioutil.WriteFile("./pkg/jsondata/metric_table.json", []byte(mockJSON), 0644)
	defer func() {
		_ = ioutil.Remove("./pkg/jsondata/metric_table.json")
	}()

	// Call the function
	repo := NewMetricTableRepositoryDb(gormDB)

	// Verify the results
	assert.NotNil(t, repo)
	assert.Equal(t, gormDB, repo.getGorm)
	assert.Equal(t, "testID", repo.metricTable.ID)
	assert.Equal(t, "testData", repo.metricTable.Data)
}