package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"sync"
	"time"

	"github.com/RedisBloom/redisbloom-go"
	"github.com/go-redis/redis/v8"
)

const (
	// Redis connection details
	redisAddr     = "localhost:6379" // Redis server address
	redisPassword = ""                // Password (leave empty if no password)
	redisDB       = 0                 // Redis database number

	// Bloom filter name
	bloomFilterName = "optimizeKeyRedisPerformance"

	// Bloom filter parameters
	n = 100000
	p = 0.0000001
	k = 25

	// CSV file details
	csvFilePath = "records.csv" // Path to the CSV file
)

var (
	ctx = context.Background()
)

// generateKey generates a single string key from CSV row fields
func generateKey(row []string) string {
	return fmt.Sprintf("%s:%s:%s:%s", row[0], row[1], row[2], row[3])
}

func main() {
	// Connect to Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
	})

	// Create a RedisBloom client
	rb := redisbloom.NewClient(rdb)

	// Calculate the number of bits required for the Bloom filter
	m := uint64(float64(-1*n*k) / math.Log(1-p))

	// Create a Bloom filter with calculated parameters
	filter := redisbloom.NewFilterOpt("BF", bloomFilterName, m, k)
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
	var mu sync.Mutex
	successfulInsertions := 0

	// Create a wait group for concurrent insertion
	var wg sync.WaitGroup

	// Process each row in the CSV file
	for {
		row, err := reader.Read()
		if err != nil {
			if err == csv.ErrFieldCount {
				break // Ignore extra fields, if any
			} else if err == io.EOF {
				break // Reached end of file
			} else {
				log.Printf("Error reading CSV row: %v\n", err)
				continue
			}
		}

		// Increment the wait group
		wg.Add(1)

		// Generate key using CSV row fields
		key := generateKey(row)

		// Concurrently insert key into Bloom filter
		go func(key string) {
			defer wg.Done()

			inserted, err := rb.AddMulti(ctx, filter, key)
			if err != nil {
				log.Printf("Failed to insert key into Bloom filter: %v\n", err)
				return
			}

			// Increment successful insertions count
			if inserted {
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
