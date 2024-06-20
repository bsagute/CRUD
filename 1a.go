package db_test

import (
    "errors"
    "net/url"
    "testing"

    "github.com/DATA-DOG/go-sqlmock"
    "github.com/stretchr/testify/assert"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "my_project/db"    // Replace with the actual path to your db package
    "my_project/model" // Replace with the actual path to your model package
)

func TestFindAllAppComp(t *testing.T) {
    // Create a new sqlmock database connection and a mock object
    dbMock, mock, err := sqlmock.New()
    assert.NoError(t, err)
    defer dbMock.Close()

    // Open a new GORM DB connection using the mock database
    gdb, err := gorm.Open(postgres.New(postgres.Config{
        Conn: dbMock,
    }), &gorm.Config{})
    assert.NoError(t, err)

    // Initialize the repository with the mocked GORM DB
    repo := db.AppComponentRepositoryDb{
        getGorm: gdb,
    }

    // Positive Case: Valid app_id
    t.Run("ValidAppID", func(t *testing.T) {
        values := url.Values{}
        values.Set("app_id", "valid_app_id")

        // Mock the getNewLayoutData function
        db.GetNewLayoutData = func(appID string, gormDB *gorm.DB) (model.AppComponent, error) {
            return model.AppComponent{
                ID: appID,
            }, nil
        }

        // Call the function
        result, err := repo.FindAllAppComp(&values)
        assert.NoError(t, err)
        assert.NotNil(t, result)
        assert.Equal(t, "valid_app_id", result.ID)
    })

    // Positive Case: Valid component_id
    t.Run("ValidComponentID", func(t *testing.T) {
        values := url.Values{}
        values.Set("component_id", "valid_component_id")

        // Mock the getNewLayoutData function
        db.GetNewLayoutData = func(appID string, gormDB *gorm.DB) (model.AppComponent, error) {
            return model.AppComponent{
                ID: appID,
            }, nil
        }

        // Call the function
        result, err := repo.FindAllAppComp(&values)
        assert.NoError(t, err)
        assert.NotNil(t, result)
        assert.Equal(t, "valid_component_id", result.ID)
    })

    // Negative Case: Error in getNewLayoutData
    t.Run("GetNewLayoutDataError", func(t *testing.T) {
        values := url.Values{}
        values.Set("app_id", "error_app_id")

        // Mock the getNewLayoutData function to return an error
        db.GetNewLayoutData = func(appID string, gormDB *gorm.DB) (model.AppComponent, error) {
            return model.AppComponent{}, errors.New("data retrieval error")
        }

        // Call the function
        result, err := repo.FindAllAppComp(&values)
        assert.Error(t, err)
        assert.EqualError(t, err, "data retrieval error")
        assert.Equal(t, model.AppComponent{}, result)
    })
}
