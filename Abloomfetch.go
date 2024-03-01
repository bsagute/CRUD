package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
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

	// Bloom filter parameters
	bloomFilterErrorRate = 0.0000001
	bloomFilterCapacity  = 100000

	// CSV file details
	csvFilePath   = "records.csv" // Path to the CSV file
	outputCSVPath = "output.csv"  // Path to the output CSV file
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

	// Close the Redis client when main function exits
	defer func() {
		if err := rdb.Close(); err != nil {
			log.Printf("Error closing Redis client: %v\n", err)
		}
	}()

	// Create a RedisBloom client
	rb := redis_bloom_go.NewClient(redisAddr, bloomFilterName, nil)

	// Reserve the Bloom filter with specified parameters
	log.Printf("Reserving Bloom filter with error rate %f and capacity %d...\n", bloomFilterErrorRate, bloomFilterCapacity)
	if err := rb.Reserve(bloomFilterName, bloomFilterErrorRate, bloomFilterCapacity); err != nil {
		log.Fatalf("Failed to reserve Bloom filter: %v", err)
	}
	fmt.Printf("Bloom filter reserved successfully: %s\n", bloomFilterName)

	// Open the CSV file
	log.Printf("Opening CSV file: %s...\n", csvFilePath)
	file, err := os.Open(csvFilePath)
	if err != nil {
		log.Fatalf("Failed to open CSV file: %v", err)
	}
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)

	// Read the CSV headers
	log.Println("Reading CSV headers...")
	headers, err := reader.Read()
	if err != nil {
		log.Fatalf("Failed to read CSV headers: %v", err)
	}

	// Append new headers for 'isFound' and 'timeToFind' columns
	headers = append(headers, "isFound", "timeToFind")

	// Create a new CSV writer for the output file
	log.Printf("Creating output CSV file: %s...\n", outputCSVPath)
	outputFile, err := os.Create(outputCSVPath)
	if err != nil {
		log.Fatalf("Failed to create output CSV file: %v", err)
	}
	defer outputFile.Close()

	// Write the headers to the output CSV file
	log.Println("Writing headers to output CSV file...")
	outputWriter := csv.NewWriter(outputFile)
	if err := outputWriter.Write(headers); err != nil {
		log.Fatalf("Failed to write headers to output CSV file: %v", err)
	}

	// Record the start time
	startTime := time.Now()

	// Process each row in the CSV file
	log.Println("Processing CSV rows...")
	for {
		row, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break // Reached end of file
			}
			log.Printf("Error reading CSV row: %v\n", err)
			continue
		}

		// Generate key using CSV row fields
		key := generateKey(row)

		// Record the start time for key search
		searchStartTime := time.Now()

		// Check if the key exists in the Bloom filter
		exists, err := rb.Exists(bloomFilterName, key)
		if err != nil {
			log.Printf("Failed to check key existence in Bloom filter: %v\n", err)
			continue
		}

		// Record the time taken to find the key
		timeToFind := time.Since(searchStartTime).Milliseconds()

		// Append 'isFound' and 'timeToFind' columns to the row
		row = append(row, fmt.Sprintf("%t", exists), fmt.Sprintf("%d", timeToFind))

		// Write the row to the output CSV file
		if err := outputWriter.Write(row); err != nil {
			log.Printf("Failed to write row to output CSV file: %v\n", err)
			continue
		}
	}

	// Flush the CSV writer buffer
	outputWriter.Flush()

	// Record the end time
	endTime := time.Now()

	// Calculate the total time taken
	totalTime := endTime.Sub(startTime)

	fmt.Printf("All keys processed successfully in %.2f seconds!\n", totalTime.Seconds())
	fmt.Printf("Output CSV file generated: %s\n", outputCSVPath)
}
