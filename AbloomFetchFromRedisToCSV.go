package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
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
	outputCSVPath = "output.csv" // Path to the output CSV file

	// Number of workers (goroutines) for concurrent processing
	numWorkers = 10
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

	// Close the Redis client when the main function exits
	defer func() {
		if err := rdb.Close(); err != nil {
			log.Printf("Error closing Redis client: %v\n", err)
		}
	}()

	// Create a RedisBloom client
	rb := redis_bloom_go.NewClient(redisAddr, bloomFilterName, nil)

	// Open the output CSV file for writing
	outputFile, err := os.Create(outputCSVPath)
	if err != nil {
		log.Fatalf("Failed to create output CSV file: %v", err)
	}
	defer outputFile.Close()

	// Create a new CSV writer
	outputWriter := csv.NewWriter(outputFile)

	// Write the CSV header
	if err := outputWriter.Write([]string{"Key", "TimeToFind(ms)", "KeySize(Bytes)"}); err != nil {
		log.Fatalf("Failed to write CSV header: %v", err)
	}

	// Flush the CSV writer buffer
	outputWriter.Flush()

	// Record the start time
	startTime := time.Now()

	// Create a wait group for workers
	var wg sync.WaitGroup

	// Create a channel to send keys to workers
	keysChan := make(chan string)

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for key := range keysChan {
				processKey(rb, rdb, outputWriter, key)
			}
		}()
	}

	// Iterate over each key in Redis
	iter := rdb.Scan(ctx, 0, "*", 0).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()
		keysChan <- key
	}

	// Close the keys channel to signal workers to stop
	close(keysChan)

	// Wait for all workers to finish
	wg.Wait()

	// Check for iterator errors
	if err := iter.Err(); err != nil {
		log.Fatalf("Redis scan iterator error: %v", err)
	}

	// Record the end time
	endTime := time.Now()

	// Calculate the total time taken
	totalTime := endTime.Sub(startTime)

	fmt.Printf("All keys processed successfully in %.2f seconds!\n", totalTime.Seconds())
	fmt.Printf("Output CSV file generated: %s\n", outputCSVPath)
}

// processKey checks if the key exists in the Bloom filter, fetches its size from Redis, and writes it to the output CSV file
func processKey(rb *redis_bloom_go.Client, rdb *redis.Client, outputWriter *csv.Writer, key string) {
	// Check if the key exists in the Bloom filter
	exists, err := rb.Exists(bloomFilterName, key)
	if err != nil {
		log.Printf("Failed to check key existence in Bloom filter: %v\n", err)
		return
	}

	if !exists {
		return // Key not found in Bloom filter
	}

	// Record the start time for key search
	searchStartTime := time.Now()

	// Fetch the key from Redis to get its size
	value, err := rdb.Get(ctx, key).Bytes()
	if err != nil {
		log.Printf("Failed to fetch key from Redis: %v\n", err)
		return
	}
	keySize := len(value)

	// Record the time taken to find the key
	timeToFind := time.Since(searchStartTime).Milliseconds()

	// Write the key, time taken, and key size to the output CSV file
	record := []string{key, strconv.FormatInt(timeToFind, 10), strconv.Itoa(keySize)}
	if err := outputWriter.Write(record); err != nil {
		log.Printf("Failed to write record to output CSV file: %v\n", err)
		return
	}

	// Flush the CSV writer buffer
	outputWriter.Flush()
}
