package validators

import (
	"digi-model-engine/utils/exceptions"
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// bind request data against the provided structs.
func ValidateRequest(c *gin.Context, model interface{}) {

	if err := c.ShouldBindJSON(&model); err != nil {
		validationErrors := make(map[string]string)
		if ve, ok := err.(validator.ValidationErrors); ok {
			for _, fieldError := range ve {
				fieldName := getFieldJSONTagName(model, fieldError.Field())
				validationErrors[fieldName] = getErrorMsg(fieldError.Tag(), fieldName)
			}
			exceptions.BadRequest(validationErrors)
		}
	}
}

// check validation error type to return proper error message
func getErrorMsg(tag, fieldName string) string {
	switch tag {
	case "required":
		return fmt.Sprintf("%s is a required field", fieldName)
	default:
		return "Request data format is not supported"
	}
}

// return json tag name
func getFieldJSONTagName(structType interface{}, fieldName string) string {
	t := reflect.TypeOf(structType)
	field, _ := t.Elem().FieldByName(fieldName)
	jsonTag := field.Tag.Get("json")
	jsonTagName := strings.Split(jsonTag, ",")[0]
	return jsonTagName
}
