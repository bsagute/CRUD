package main

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func convertToString(value interface{}) (string, error) {
	// Use reflection to handle various types
	switch v := value.(type) {
	case string:
		return v, nil
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v), nil
	case float32, float64:
		return strconv.FormatFloat(v.(float64), 'f', -1, 64), nil
	default:
		return "", fmt.Errorf("unsupported type: %s", reflect.TypeOf(value))
	}
}

func hashMetricsData(fields ...interface{}) (string, error) {
	// Convert all fields to strings
	var stringFields []string
	for _, field := range fields {
		strField, err := convertToString(field)
		if err != nil {
			return "", err
		}
		stringFields = append(stringFields, strField)
	}

	// Concatenate string fields
	dataToHash := strings.Join(stringFields, "")

	// Hash the concatenated data
	hashedData, err := bcrypt.GenerateFromPassword([]byte(dataToHash), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	// Convert hashed data to string
	hashedString := string(hashedData)

	return hashedString, nil
}

func main() {
	// Example input values
	entityID := "exampleEntity"
	metricsID := "exampleMetrics"
	metricTimestamp := "2023-01-01T12:00:00Z"
	metricValue := 3.14

	// Call the hashMetricsData function with different types of parameters
	hashedResult, err := hashMetricsData(entityID, metricsID, metricTimestamp, metricValue)
	if err != nil {
		log.Fatal(err)
	}

	// Display the hashed result
	fmt.Println("Hashed Result:", hashedResult)
}
