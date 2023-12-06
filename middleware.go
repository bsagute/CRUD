package middlewares

import (
	"digi-model-engine/models"
	"digi-model-engine/utils/exceptions"
	"digi-model-engine/utils/logger"

	"net/http"

	"github.com/gin-gonic/gin"
)

// This function is a middleware, it recover all the panic causes in all over the project and notify bugsnag if there are any error
func ExceptionHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		// this will be called after all the code completed or the panic raised
		defer func() {
			if r := recover(); r != nil {

				// Initializing the defualt error status code to sent in response
				// getting the error string which is sent via panic
				errCode := http.StatusInternalServerError
				var err error
				var errMsg interface{}

				if ce, ok := r.(exceptions.CustomError); ok {
					errCode = ce.StatusCode
					err = ce.Err
					errMsg = ce.ErrorMsg
				}

				// Error response model is creation with the values provided
				resp := models.Response{
					Success:      false,
					Message:      "fail",
					ResponseCode: errCode,
					Data:         nil,
					Error:        errMsg,
				}

				logger.Error(err, errMsg)
				// Sending error and abort response
				c.AbortWithStatusJSON(errCode, resp)
			}
		}()
		c.Next()
	}
}
