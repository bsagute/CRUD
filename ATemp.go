package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize Gin router
	router := gin.Default()

	// Define route handler
	router.GET("/:tableName", handleRequest)

	// Start server
	fmt.Println("Server started at localhost:9090")
	router.Run(":9090")
}

func handleRequest(c *gin.Context) {
	// Parse query parameters
	tableName := c.Param("tableName")
	queryParams := c.Request.URL.Query()

	// Validate tableName parameter
	if tableName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tableName parameter is missing"})
		return
	}

	// Generate SQL query and Redis key
	query, key := generateQueryAndKey(tableName, queryParams)

	// Simulate fetching data from PostgreSQL (dummy call)
	dataFromPostgres, err := fetchDataFromPostgres(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching data from PostgreSQL"})
		return
	}

	// Convert response to JSON
	jsonData, err := json.Marshal(dataFromPostgres)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error converting data to JSON"})
		return
	}

	// Simulate storing JSON data in Redis (dummy call)
	err = storeJSONInRedis(key, jsonData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error storing JSON data in Redis"})
		return
	}

	// Return JSON response
	c.JSON(http.StatusOK, gin.H{"data": string(jsonData)})
}

func generateQueryAndKey(tableName string, queryParams map[string][]string) (string, string) {
	// Prepare SQL query
	query := "SELECT * FROM " + tableName

	// Prepare Redis key
	key := fmt.Sprintf("%s?", tableName)
	for key, values := range queryParams {
		for _, value := range values {
			key += fmt.Sprintf("%s=%s&", key, value)
		}
	}
	// Remove the trailing "&"
	key = key[:len(key)-1]

	return query, key
}

func fetchDataFromPostgres(query string) (interface{}, error) {
	// Placeholder for fetching data from PostgreSQL (dummy call)
	// In a real-world scenario, this function should execute the SQL query on PostgreSQL and return the result
	return map[string]interface{}{
		"example_field": "example_value",
	}, nil
}

func storeJSONInRedis(key string, jsonData []byte) error {
	// Placeholder for storing JSON data in Redis (dummy call)
	// In a real-world scenario, this function should store the JSON data in Redis
	fmt.Printf("Storing JSON data with key '%s' in Redis...\n", key)
	return nil
}
