package exceptions

import (
	"errors"
	"fmt"
)

// This structs extends the available error variable. It is used to represent all the panic and raise a specific format error
type CustomError struct {
	Err        error
	ErrorMsg   interface{}
	StatusCode int
}

func (ce CustomError) Error() string {
	return fmt.Sprintf("%s: %s ", ce.Err, ce.ErrorMsg)
}

// This function accepts any type of error as parameter and raise error with StatusCode 400
func BadRequest(errMsg interface{}) {
	raiseError(errors.New("Bad Request"), errMsg, 400)
}

// This function accepts any type of error as parameter and raise error with StatusCode 404
func NotFoundError(errMsg interface{}) {
	raiseError(errors.New("Not Found Error"), errMsg, 404)
}

// This function accepts any type of error as parameter and raise error with StatusCode 404
func InternalServerError(errMsg interface{}) {
	raiseError(errors.New("Internal Server Error"), errMsg, 500)
}

// This function accepts any type of error and StatusCode as parameter and raise panic with customError model
func raiseError(err error, errMsg interface{}, statusCode int) {
	ce := CustomError{
		Err:        err,
		ErrorMsg:   errMsg,
		StatusCode: statusCode,
	}
	panic(ce)
}
