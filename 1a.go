// main_test.go
package main

import (
    "testing"
)

func TestCheckError(t *testing.T) {
    t.Run("should panic if error is not nil", func(t *testing.T) {
        defer func() {
            if r := recover(); r == nil {
                t.Errorf("Expected panic but code did not panic")
            }
        }()
        CheckError(fmt.Errorf("this is an error"))
    })

    t.Run("should not panic if error is nil", func(t *testing.T) {
        defer func() {
            if r := recover(); r != nil {
                t.Errorf("Did not expect panic but code panicked")
            }
        }()
        CheckError(nil)
    })
}
