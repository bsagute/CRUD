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

	// Bloom filter name
	bloomFilterName = "optimizeKeyRedisPerformance"
)

var (
	RedisClient *redis.Client
	ctx         = context.Background()
)

func main() {
	// Initialize Redis connection
	if err := initRedisDB(); err != nil {
		log.Fatalf("Failed to initialize Redis DB: %v", err)
	}

	// Create a new Bloom filter
	if err := createBloomFilter(); err != nil {
		log.Fatalf("Failed to create Bloom filter: %v", err)
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
	if _, err := reader.Read(); err != nil {
		log.Fatalf("Failed to read CSV headers: %v", err)
	}

	// Record the start time
	startTime := time.Now()

	// Record the number of successful insertions
	successfulInsertions := 0

	// Process each row in the CSV file
	var wg sync.WaitGroup
	for {
		row, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("Error reading CSV row: %v", err)
		}

		wg.Add(1)
		go func(row []string) {
			defer wg.Done()

			// Construct MetricData object from CSV row
			metricData := constructMetricData(row)

			// Generate key from MetricData fields
			key := generateKey(metricData)

			// Check if the key exists in the Bloom filter
			if exists, err := existsInBloomFilter(key); err != nil {
				log.Printf("Error checking key in Bloom filter: %v", err)
			} else if !exists {
				// If the key does not exist in the Bloom filter, add it and insert into Redis
				if err := insertIntoBloomFilter(key); err != nil {
					log.Printf("Failed to insert key into Bloom filter: %v", err)
				} else {
					if err := insertIntoRedis(key); err != nil {
						log.Printf("Failed to insert key into Redis: %v", err)
					} else {
						successfulInsertions++
					}
				}
			}
		}(row)
	}

	wg.Wait()

	// Record the end time
	endTime := time.Now()

	// Calculate the total time taken
	totalTime := endTime.Sub(startTime)

	fmt.Printf("All keys processed successfully in %v seconds!\n", totalTime.Seconds())
	fmt.Printf("Number of successful insertions: %d\n", successfulInsertions)
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

// createBloomFilter creates a new Bloom filter if not exists
func createBloomFilter() error {
	if RedisClient == nil {
		return fmt.Errorf("Redis client is not initialized")
	}

	// Check if the Bloom filter already exists
	exists, err := RedisClient.Do(ctx, "BF.EXISTS", bloomFilterName).Bool()
	if err != nil {
		return fmt.Errorf("Error checking Bloom filter existence: %v", err)
	}

	// If the Bloom filter doesn't exist, create a new one
	if !exists {
		if _, err := RedisClient.Do(ctx, "BF.RESERVE", bloomFilterName, "0.001", "1000000").Result(); err != nil {
			return fmt.Errorf("Error creating Bloom filter: %v", err)
		}
	}

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

// insertIntoBloomFilter inserts a key into the Bloom filter
func insertIntoBloomFilter(key string) error {
	if RedisClient == nil {
		return fmt.Errorf("Redis client is not initialized")
	}

	// Insert the key into the Bloom filter
	if _, err := RedisClient.Do(ctx, "BF.ADD", bloomFilterName, key).Result(); err != nil {
		return fmt.Errorf("Error inserting key into Bloom filter: %v", err)
	}

	return nil
}

// existsInBloomFilter checks if a key exists in the Bloom filter
func existsInBloomFilter(key string) (bool, error) {
	if RedisClient == nil {
		return false, fmt.Errorf("Redis client is not initialized")
	}

	// Check if the key exists in the Bloom filter
	exists, err := RedisClient.Do(ctx, "BF.EXISTS", bloomFilterName, key).Int()
	if err != nil {
		return false, fmt.Errorf("Error checking key in Bloom filter: %v", err)
	}

	return exists == 1, nil
}

// insertIntoRedis inserts a key into Redis
func insertIntoRedis(key string) error {
	if RedisClient == nil {
		return fmt.Errorf("Redis client is not initialized")
	}

	// Insert the key into Redis
	if _, err := RedisClient.Do(ctx, "SET", key, "").Result(); err != nil {
		return fmt.Errorf("Error inserting key into Redis: %v", err)
	}

	return nil
}
