package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/RedisBloom/redisbloom-go"
	"github.com/go-redis/redis/v8"
)

const (
	// CSV file details
	csvFilePath = "records_1000.csv" // Path to the CSV file

	// Redis connection details
	redisAddr     = "localhost:6379" // Redis server address
	redisPassword = ""                // Password (leave empty if no password)
	redisDB       = 0                 // Redis database number

	// Bloom filter name
	bloomFilterName = "optimizeKeyRedisPerformance"
)

var (
	ctx = context.Background()
)

func main() {
	// Connect to Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
	})

	// Create a RedisBloom client
	rb := redisbloom.NewClient(rdb)

	// Create a Bloom filter
	filter := redisbloom.NewFilterOpt("BF", bloomFilterName, 1000000, 0.001)
	if err := rb.CreateFilter(ctx, filter); err != nil {
		log.Fatalf("Failed to create Bloom filter: %v", err)
	}
	fmt.Printf("Bloom filter created successfully: %s\n", bloomFilterName)

	// Open the CSV file
	file, err := os.Open(csvFilePath)
	if err != nil {
		log.Fatalf("Failed to open CSV file: %v", err)
	}
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)

	// Read the CSV headers
	if _, err := reader.Read(); err != nil {
		log.Fatalf("Failed to read CSV headers: %v", err)
	}

	// Record the start time
	startTime := time.Now()

	// Record the number of successful insertions
	successfulInsertions := 0

	// Process each row in the CSV file
	for {
		row, err := reader.Read()
		if err != nil {
			if err == csv.ErrFieldCount || err == io.EOF {
				break
			}
			log.Fatalf("Error reading CSV row: %v", err)
		}

		// Construct MetricData object from CSV row
		metricData := constructMetricData(row)

		// Generate key from MetricData fields
		key := generateKey(metricData)

		// Insert key into Bloom filter
		if inserted, err := rb.AddMulti(ctx, filter, key); err != nil {
			log.Printf("Failed to insert key into Bloom filter: %v", err)
		} else if inserted {
			// If the key was added to the Bloom filter, insert it into Redis
			if err := rdb.Set(ctx, key, "", 0).Err(); err != nil {
				log.Printf("Failed to insert key into Redis: %v", err)
			} else {
				successfulInsertions++
			}
		}
	}

	// Record the end time
	endTime := time.Now()

	// Calculate the total time taken
	totalTime := endTime.Sub(startTime)

	fmt.Printf("All keys processed successfully in %v seconds!\n", totalTime.Seconds())
	fmt.Printf("Number of successful insertions: %d\n", successfulInsertions)
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
