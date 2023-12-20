package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type MetricData struct {
	MetricId        string `json:"metricId"`
	EntityId        string `json:"entityId"`
	MetricValue     float64   `json:"metricValue"`
	MetricTimestamp string `json:"metricTimestamp"`
}

func main() {
	// Replace these values with your actual data
	data := MetricData{
		MetricId:        "your_metric_id",
		EntityId:        "your_entity_id",
		MetricValue:     123.45,
		MetricTimestamp: time.Now().Format("2006-01-02T15:04:05.999999999Z07:00"), // Format as RFC3339
	}

	// Convert struct to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	// Define the URL of the API you want to send the data to
	apiURL := "http://localhost:8080/api/insert-metric-data-redis-db"

	// Make the HTTP POST request
	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error making HTTP request:", err)
		return
	}
	defer resp.Body.Close()

	// Print the response status
	fmt.Println("Response Status:", resp.Status)
}
