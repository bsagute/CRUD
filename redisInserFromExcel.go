package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	// CSV file details
	csvFilePath = "records_1000.csv" // Path to the CSV file

	// Redis connection details
	redisAddr     = "localhost:6379" // Redis server address
	redisPassword = ""               // Password (leave empty if no password)
	redisDB       = 0                // Redis database number
)

type MetricData struct {
	EntityID    string  `json:"entityId"`
	MetricValue float64 `json:"metricValue"`
	MetricID    string  `json:"metricId"`
	Timestamp   string  `json:"timestamp"`
}

var (
	RedisClient *redis.Client
	ctx         = context.Background()
)

func main() {
	// Initialize Redis connection
	if err := initRedisDB(); err != nil {
		log.Fatalf("Failed to initialize Redis DB: %v", err)
	}

	// Open the CSV file
	file, err := os.Open(csvFilePath)
	if err != nil {
		log.Fatalf("Failed to open CSV file: %v", err)
	}
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)

	// Read the CSV headers
	_, err = reader.Read()
	if err != nil {
		log.Fatalf("Failed to read CSV headers: %v", err)
	}

	// Create a channel to communicate between goroutines
	keyCh := make(chan string)

	// Create a wait group to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Record the start time
	startTime := time.Now()

	// Record the number of failed insertions
	failedInsertions := 0

	// Process each row in the CSV file concurrently
	for rowCount := 1; ; rowCount++ {
		// Read the next row from the CSV file
		row, err := reader.Read()
		if err != nil {
			// Check for end of file
			if err == io.EOF {
				break
			}
			log.Fatalf("Error reading CSV row: %v", err)
		}

		// Increment the wait group counter
		wg.Add(1)

		// Process the row in a separate goroutine
		go func(row []string, rowCount int) {
			defer wg.Done()

			// Construct MetricData object from CSV row
			metricData := constructMetricData(row)

			// Generate key from MetricData fields
			key := generateKey(metricData)

			// Send the key through the channel
			keyCh <- key
		}(row, rowCount)
	}

	// Close the channel after all keys are sent
	go func() {
		wg.Wait()
		close(keyCh)
	}()

	// Insert keys into Redis from the channel
	for key := range keyCh {
		if err := insertKeyIntoRedis(key); err != nil {
			log.Printf("Failed to insert key into Redis: %v", err)
			failedInsertions++
		}
	}

	// Record the end time
	endTime := time.Now()

	// Calculate the total time taken
	totalTime := endTime.Sub(startTime)

	fmt.Printf("All keys inserted into Redis successfully in %v seconds!\n", totalTime.Seconds())
	fmt.Printf("Number of failed insertions: %d\n", failedInsertions)
}

// InitializeRedisDB initializes the Redis database connection
func initRedisDB() error {
	// Create a Redis client
	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
	})

	// Ping the Redis server to test the connection
	if err := client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("Failed to connect to Redis: %v", err)
	}

	// Assign the Redis client to the global variable for later use
	RedisClient = client

	return nil
}

// constructMetricData constructs a MetricData object from a CSV row
func constructMetricData(row []string) MetricData {
	entityID := strings.Join(row, ":")
	metricValue := 0.0               // Replace with actual value conversion from row
	metricID := ""                   // Replace with actual value conversion from row
	timestamp := time.Now().String() // Convert timestamp to string

	return MetricData{
		EntityID:    entityID,
		MetricValue: metricValue,
		MetricID:    metricID,
		Timestamp:   timestamp,
	}
}

// generateKey generates a single string key from MetricData fields
func generateKey(data MetricData) string {
	return fmt.Sprintf("%s:%f:%s:%s", data.EntityID, data.MetricValue, data.MetricID, data.Timestamp)
}

// insertKeyIntoRedis inserts a key into Redis using HSet
func insertKeyIntoRedis(key string) error {
	if RedisClient == nil {
		return fmt.Errorf("Redis client is not initialized")
	}

	// Insert the key into Redis using HSet
	if _, err := RedisClient.HSet(ctx, "metric_data", key, "").Result(); err != nil {
		return fmt.Errorf("Error inserting key into Redis: %v", err)
	}

	return nil
}
