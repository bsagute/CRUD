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

	redis_bloom_go "github.com/RedisBloom/redisbloom-go"
	"github.com/go-redis/redis/v8"
)

const (
	// Redis connection details
	redisAddr     = "localhost:6379" // Redis server address
	redisPassword = ""               // Password (leave empty if no password)
	redisDB       = 9                // Redis database number

	// Bloom filter name
	bloomFilterName = "optimizeKeyRedisPerformance"

	// CSV file details
	csvFilePath = "records.csv" // Path to the CSV file
)

var (
	ctx = context.Background()
)

// generateKey generates a single string key from CSV row fields
func generateKey(row []string) string {
	return strings.Join(row, ":")
}

func main() {
	// Connect to Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
	})

	// Close the Redis client when main function exits
	defer func() {
		if err := rdb.Close(); err != nil {
			log.Printf("Error closing Redis client: %v\n", err)
		}
	}()

	// Create a RedisBloom client
	rb := redis_bloom_go.NewClient(redisAddr, bloomFilterName, nil)

	// Reserve the Bloom filter with specified parameters
	if err := rb.Reserve(bloomFilterName, 0.0000001, 100000); err != nil {
		log.Fatalf("Failed to reserve Bloom filter: %v", err)
	}
	fmt.Printf("Bloom filter reserved successfully: %s\n", bloomFilterName)

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
	var mu sync.Mutex
	successfulInsertions := 0

	// Create a wait group for concurrent insertion
	var wg sync.WaitGroup

	// Process each row in the CSV file
	for {
		row, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break // Reached end of file
			}
			log.Printf("Error reading CSV row: %v\n", err)
			continue
		}

		// Increment the wait group
		wg.Add(1)

		// Generate key using CSV row fields
		key := generateKey(row)

		// Concurrently insert key into Bloom filter
		go func(key string) {
			defer wg.Done()

			exists, err := rb.Add(key, "") // Use key as the only argument for the Add method
			if err != nil {
				log.Printf("Failed to insert key into Bloom filter: %v\n", err)
				return
			}

			// Increment successful insertions count
			if !exists {
				mu.Lock()
				defer mu.Unlock()
				successfulInsertions++
			}
		}(key)
	}

	// Wait for all insertions to complete
	wg.Wait()

	// Record the end time
	endTime := time.Now()

	// Calculate the total time taken
	totalTime := endTime.Sub(startTime)

	fmt.Printf("All keys processed successfully in %.2f seconds!\n", totalTime.Seconds())
	fmt.Printf("Number of successful insertions: %d\n", successfulInsertions)
}
