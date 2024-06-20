package db_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/google/uuid"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "my_project/db"    // Replace with the actual path to your db package
    "my_project/model" // Replace with the actual path to your model package
    "github.com/DATA-DOG/go-sqlmock"
)

func TestGetBreadCumb(t *testing.T) {
    // Create a new sqlmock database connection and a mock object
    dbMock, mock, err := sqlmock.New()
    assert.NoError(t, err)
    defer dbMock.Close()

    // Open a new GORM DB connection using the mock database
    gdb, err := gorm.Open(postgres.New(postgres.Config{
        Conn: dbMock,
    }), &gorm.Config{})
    assert.NoError(t, err)

    // Initialize navigation slice
    navigation := []model.EntityInfo{}

    // Test case: Entity type is "framework"
    t.Run("EntityTypeFramework", func(t *testing.T) {
        entityID := uuid.New()
        mockEntityData := model.EntityInfo{
            EntityId:   entityID,
            EntityType: "framework",
        }

        // Mock the database query
        mock.ExpectQuery(`SELECT \* FROM "entity_info" WHERE "entity_info"."entity_id" = \$1 ORDER BY "entity_info"."entity_id" LIMIT 1`).
            WithArgs(entityID).
            WillReturnRows(sqlmock.NewRows([]string{"entity_id", "entity_type"}).
                AddRow(mockEntityData.EntityId, mockEntityData.EntityType))

        // Call the function
        resultNavigation, resultLen := db.GetBreadCumb(entityID, navigation, gdb)

        // Verify the results
        assert.Equal(t, 1, resultLen)
        assert.Equal(t, mockEntityData, resultNavigation[0])
    })

    // Test case: Entity type is not "framework"
    t.Run("EntityTypeNotFramework", func(t *testing.T) {
        entityID := uuid.New()
        mockEntityData := model.EntityInfo{
            EntityId:      entityID,
            EntityType:    "other",
            ParentEntityId: uuid.New(),
        }

        // Mock the database query
        mock.ExpectQuery(`SELECT \* FROM "entity_info" WHERE "entity_info"."entity_id" = \$1 ORDER BY "entity_info"."entity_id" LIMIT 1`).
            WithArgs(entityID).
            WillReturnRows(sqlmock.NewRows([]string{"entity_id", "entity_type", "parent_entity_id"}).
                AddRow(mockEntityData.EntityId, mockEntityData.EntityType, mockEntityData.ParentEntityId))

        // Mock the recursive call
        db.GetBreadCumb = func(entityID uuid.UUID, navigation []model.EntityInfo, getGorm *gorm.DB) ([]model.EntityInfo, int) {
            return append([]model.EntityInfo{mockEntityData}, navigation...), 1
        }

        // Call the function
        resultNavigation, resultLen := db.GetBreadCumb(entityID, navigation, gdb)

        // Verify the results
        assert.Equal(t, 1, resultLen)
        assert.Equal(t, mockEntityData, resultNavigation[0])
    })

    // Test case: Error in database query
    t.Run("DatabaseQueryError", func(t *testing.T) {
        entityID := uuid.New()

        // Mock the database query to return an error
        mock.ExpectQuery(`SELECT \* FROM "entity_info" WHERE "entity_info"."entity_id" = \$1 ORDER BY "entity_info"."entity_id" LIMIT 1`).
            WithArgs(entityID).
            WillReturnError(errors.New("database query error"))

        // Call the function
        resultNavigation, resultLen := db.GetBreadCumb(entityID, navigation, gdb)

        // Verify the results
        assert.Equal(t, 0, resultLen)
        assert.Empty(t, resultNavigation)
    })
}
