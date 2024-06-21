package yourpackage

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// Define the model.AppComponent struct if not already defined
type AppComponent struct {
	SomeField string
}

// Mock function to simulate getComponentList
func getComponentListMock(app_id string) (AppComponent, error) {
	return AppComponent{SomeField: "expected_value"}, nil
}

// Mock function to simulate getNewLayoutData
func getNewLayoutDataMock(app_id string, db *gorm.DB) (AppComponent, error) {
	return AppComponent{SomeField: "expected_value"}, nil
}

// Mock the AppComponentRepositoryDb struct
type AppComponentRepositoryDb struct {
	appinfo AppComponent
	getGorm *gorm.DB
}

// Mock the FindAllAppComp method
func (s *AppComponentRepositoryDb) FindAllAppComp(level *url.Values) (AppComponent, error) {
	var err error
	var payload AppComponent

	app_id := level.Get("app_id")

	if app_id == "" {
		app_id = level.Get("component_id")
	}

	// payload, err := getComponentList(app_id)
	// payload, err = getLayoutData(app_id, s.getGorm)
	payload, err = getNewLayoutData(app_id, s.getGorm)

	if err != nil {
		return payload, err
	}

	return payload, nil
}

// Unit test for FindAllAppComp function
func TestFindAllAppComp(t *testing.T) {
	repo := &AppComponentRepositoryDb{
		// Mock the db connection if necessary
	}

	// Mock the functions
	getComponentList = getComponentListMock
	getNewLayoutData = getNewLayoutDataMock

	level := url.Values{}
	level.Set("app_id", "test_app_id")

	payload, err := repo.FindAllAppComp(&level)

	// Assert the results
	assert.Nil(t, err)
	assert.NotNil(t, payload)
	assert.Equal(t, "expected_value", payload.SomeField) // Change "SomeField" to actual field name
}
