package db_test

import (
    "testing"

    "github.com/stretchr/testify/assert"
    "my_project/model" // Replace with the actual path to your model package
    "my_project/db"    // Replace with the actual path to your db package
)

func TestRemove(t *testing.T) {
    // Test case: Remove an element from the middle of the slice
    t.Run("RemoveMiddleElement", func(t *testing.T) {
        input := []model.AppComponentLayout{
            {ID: "1"}, {ID: "2"}, {ID: "3"}, {ID: "4"}, {ID: "5"},
        }
        expected := []model.AppComponentLayout{
            {ID: "1"}, {ID: "2"}, {ID: "5"}, {ID: "4"},
        }

        result := db.Remove(input, 2)
        assert.Equal(t, expected, result)
    })

    // Test case: Remove the first element of the slice
    t.Run("RemoveFirstElement", func(t *testing.T) {
        input := []model.AppComponentLayout{
            {ID: "1"}, {ID: "2"}, {ID: "3"},
        }
        expected := []model.AppComponentLayout{
            {ID: "3"}, {ID: "2"},
        }

        result := db.Remove(input, 0)
        assert.Equal(t, expected, result)
    })

    // Test case: Remove the last element of the slice
    t.Run("RemoveLastElement", func(t *testing.T) {
        input := []model.AppComponentLayout{
            {ID: "1"}, {ID: "2"}, {ID: "3"},
        }
        expected := []model.AppComponentLayout{
            {ID: "1"}, {ID: "2"},
        }

        result := db.Remove(input, 2)
        assert.Equal(t, expected, result)
    })

    // Test case: Index out of bounds (should not alter the slice)
    t.Run("IndexOutOfBounds", func(t *testing.T) {
        input := []model.AppComponentLayout{
            {ID: "1"}, {ID: "2"}, {ID: "3"},
        }
        expected := []model.AppComponentLayout{
            {ID: "1"}, {ID: "2"}, {ID: "3"},
        }

        result := db.Remove(input, 5)
        assert.Equal(t, expected, result)
    })
}
