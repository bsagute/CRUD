package main

import (
	"fmt"
	"log"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

func hashMetricsData(entityID, metricsID, metricTimestamp string, metricValue float64) (string, error) {
	// Convert float to string
	metricValueStr := strconv.FormatFloat(metricValue, 'f', -1, 64)

	// Concatenate input values
	dataToHash := entityID + metricsID + metricTimestamp + metricValueStr

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

	// Call the hashMetricsData function
	hashedResult, err := hashMetricsData(entityID, metricsID, metricTimestamp, metricValue)
	if err != nil {
		log.Fatal(err)
	}

	// Display the hashed result
	fmt.Println("Hashed Result:", hashedResult)
}
