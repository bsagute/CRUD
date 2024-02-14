package main

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"unsafe"
)

type Data struct {
	RedisKy     string  `json:"redisKy"`
	MetricValue float64 `json:"metricvalue"`
	MetricID    string  `json:"metricId"`
	Timestamp   string  `json:"timestamp"`
}

func main() {
	// Example data
	exampleData := Data{
		RedisKy:     "28cf71a2-3da0-4fd3-87e3-af8e1b0f7b24",
		MetricValue: 30.0,
		MetricID:    "28cf71a2-3da0-4fd3-87e3-af8e1b0f7b24",
		Timestamp:   "2024-02-12T12:00:00Z",
	}

	// Convert example data to JSON
	jsonData, err := json.Marshal(exampleData)
	if err != nil {
		log.Fatalf("Error marshalling JSON: %v", err)
	}

	// Compress JSON data
	var compressedData bytes.Buffer
	gzipWriter := gzip.NewWriter(&compressedData)
	if _, err := gzipWriter.Write(jsonData); err != nil {
		log.Fatalf("Error compressing data: %v", err)
	}
	if err := gzipWriter.Close(); err != nil {
		log.Fatalf("Error closing gzip writer: %v", err)
	}

	// Calculate size of compressed JSON data
	dataSize := len(compressedData.Bytes())
	fmt.Printf("Size of compressed JSON data: %d bytes\n", dataSize)

	// Calculate size of string keys
	redisKeySize := int(unsafe.Sizeof(exampleData.RedisKy))    // Convert to int
	metricIDSize := int(unsafe.Sizeof(exampleData.MetricID))   // Convert to int
	timestampSize := int(unsafe.Sizeof(exampleData.Timestamp)) // Convert to int

	// Calculate size of SHA hash key
	hash := sha256.Sum256([]byte(exampleData.RedisKy))
	shaKeySize := int(unsafe.Sizeof(hash)) // Convert to int

	fmt.Printf("Size of Redis Key: %d bytes\n", redisKeySize)
	fmt.Printf("Size of Metric ID: %d bytes\n", metricIDSize)
	fmt.Printf("Size of Timestamp: %d bytes\n", timestampSize)
	fmt.Printf("Size of SHA256 Key: %d bytes\n", shaKeySize)

	// Calculate total size
	totalKeySize := redisKeySize + metricIDSize + timestampSize + shaKeySize
	totalSize := totalKeySize + dataSize
	fmt.Printf("Total size: %d bytes\n", totalSize)

	// Print addition of all struct fields
	addition := redisKeySize + metricIDSize + timestampSize + int(unsafe.Sizeof(exampleData.MetricValue))
	fmt.Printf("Addition of all struct fields: %d bytes\n", addition)

	// Calculate SHA hash of combined data
	combinedData := exampleData.RedisKy + exampleData.MetricID + exampleData.Timestamp
	combinedHash := sha256.Sum256([]byte(combinedData))
	fmt.Printf("SHA256 of combined data: %x\n", combinedHash)

	// Calculate size of SHA hash of combined data
	combinedHashSize := int(unsafe.Sizeof(combinedHash)) // Convert to int
	fmt.Printf("Size of SHA256 of combined data: %d bytes\n", combinedHashSize)

	// Example hash key in Redis
	hashKey := hex.EncodeToString(hash[:])
	fmt.Printf("Example Hash Key in Redis: %s\n", hashKey)

	// Example Redis Bloom Filter
	bloomFilterKey := "example_bloom_filter_key" // Example Bloom Filter key
	bloomFilterSize := 100                       // Example Bloom Filter size
	fmt.Printf("Size of Bloom Filter Key in Redis: %d bytes\n", len(bloomFilterKey))
	fmt.Printf("Size of Bloom Filter Object: %d bytes\n", bloomFilterSize)
}
