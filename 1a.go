package db_test

import (
    "encoding/json"
    "io/ioutil"
    "os"
    "testing"

    "github.com/stretchr/testify/assert"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "my_project/db"    // Replace with the actual path to your db package
    "my_project/model" // Replace with the actual path to your model package
)

func TestNewAppComponentRepositoryDbLevelOne(t *testing.T) {
    // Create a temporary JSON file with mock data
    mockData := `{
        "ID": "test_id",
        "Name": "Test Component"
    }`
    tmpFile, err := ioutil.TempFile(os.TempDir(), "test_*.json")
    assert.NoError(t, err)
    defer os.Remove(tmpFile.Name())

    _, err = tmpFile.Write([]byte(mockData))
    assert.NoError(t, err)
    tmpFile.Close()

    // Modify the function to read from the temporary file
    originalFilePath := "./pkg/jsondata/app_component_details_level_0.json"
    newFilePath := tmpFile.Name()

    os.Rename(originalFilePath, originalFilePath+".bak") // Backup original file
    defer os.Rename(originalFilePath+".bak", originalFilePath) // Restore original file after test
    os.Rename(newFilePath, originalFilePath)

    // Open a new GORM DB connection using a mocked database
    dbMock, _, err := sqlmock.New()
    assert.NoError(t, err)
    defer dbMock.Close()

    gdb, err := gorm.Open(postgres.New(postgres.Config{
        Conn: dbMock,
    }), &gorm.Config{})
    assert.NoError(t, err)

    // Call the function
    result := db.NewAppComponentRepositoryDbLevelOne(gdb)

    // Verify the results
    assert.Equal(t, "test_id", result.appinfo.ID)
    assert.Equal(t, "Test Component", result.appinfo.Name)
    assert.NotNil(t, result.getGorm)
}

func TestNewAppComponentRepositoryDbLevelOne_FileError(t *testing.T) {
    // Open a new GORM DB connection using a mocked database
    dbMock, _, err := sqlmock.New()
    assert.NoError(t, err)
    defer dbMock.Close()

    gdb, err := gorm.Open(postgres.New(postgres.Config{
        Conn: dbMock,
    }), &gorm.Config{})
    assert.NoError(t, err)

    // Temporarily rename the JSON file to simulate file read error
    originalFilePath := "./pkg/jsondata/app_component_details_level_0.json"
    os.Rename(originalFilePath, originalFilePath+".bak") // Backup original file
    defer os.Rename(originalFilePath+".bak", originalFilePath) // Restore original file after test

    assert.PanicsWithError(t, "json file error open ./pkg/jsondata/app_component_details_level_0.json: no such file or directory", func() {
        db.NewAppComponentRepositoryDbLevelOne(gdb)
    })
}

func TestNewAppComponentRepositoryDbLevelOne_UnmarshalError(t *testing.T) {
    // Create a temporary JSON file with invalid data
    invalidData := `{invalid_json}`
    tmpFile, err := ioutil.TempFile(os.TempDir(), "test_invalid_*.json")
    assert.NoError(t, err)
    defer os.Remove(tmpFile.Name())

    _, err = tmpFile.Write([]byte(invalidData))
    assert.NoError(t, err)
    tmpFile.Close()

    // Modify the function to read from the temporary file
    originalFilePath := "./pkg/jsondata/app_component_details_level_0.json"
    newFilePath := tmpFile.Name()

    os.Rename(originalFilePath, originalFilePath+".bak") // Backup original file
    defer os.Rename(originalFilePath+".bak", originalFilePath) // Restore original file after test
    os.Rename(newFilePath, originalFilePath)

    // Open a new GORM DB connection using a mocked database
    dbMock, _, err := sqlmock.New()
    assert.NoError(t, err)
    defer dbMock.Close()

    gdb, err := gorm.Open(postgres.New(postgres.Config{
        Conn: dbMock,
    }), &gorm.Config{})
    assert.NoError(t, err)

    assert.PanicsWithError(t, "json unmarshal error invalid character 'i' looking for beginning of object key string", func() {
        db.NewAppComponentRepositoryDbLevelOne(gdb)
    })
}
