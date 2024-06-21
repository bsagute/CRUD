package db_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"slices"
	"strings"
)

// Mock the model.Applist struct if not already defined
type Applist struct {
	// Define fields according to your actual model
	SomeField string
}

// Mock function to simulate getApplistData
func getApplistDataMock(db *gorm.DB) (Applist, error) {
	return Applist{SomeField: "expected_value"}, nil
}

// Mock the ApplistRepositoryDb struct
type ApplistRepositoryDb struct {
	appinfo Applist
	getGorm *gorm.DB
}

// Implement the FindAllComp method
func (s *ApplistRepositoryDb) FindAllComp() (Applist, error) {
	var entityList Applist
	entityList, err := getApplistData(s.getGorm)

	if err != nil {
		log.Printf("Error %s", err.Error())
		return entityList, err
	}

	return entityList, nil
}

// Unit test for FindAllComp function
func TestFindAllComp(t *testing.T) {
	repo := &ApplistRepositoryDb{
		// Mock the db connection if necessary
	}

	// Mock the function
	getApplistData = getApplistDataMock

	entityList, err := repo.FindAllComp()

	// Assert the results
	assert.Nil(t, err)
	assert.NotNil(t, entityList)
	assert.Equal(t, "expected_value", entityList.SomeField) // Change "SomeField" to actual field name
}
