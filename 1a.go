package db_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "my_project/db" // Replace with the actual path to your db package
)

func TestPrettyPrint(t *testing.T) {
    t.Run("ValidInput", func(t *testing.T) {
        input := map[string]string{"key": "value"}
        err := db.PrettyPrint(input)
        assert.NoError(t, err)
    })

    t.Run("InvalidInput", func(t *testing.T) {
        // Intentionally passing an invalid input to cause JSON marshaling error
        input := make(chan int)
        err := db.PrettyPrint(input)
        assert.Error(t, err)
    })
}
