package main

import (
	"bytes"
	"compress/zlib"
	"context"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/go-redis/redis/v8"
)

// Emp struct represents your employee structure
type Emp struct {
	Id         int
	IsExist    bool
	Name       string
	CreatedAt  time.Time
	FloatValue float64
}

func compress(emp *Emp) ([]byte, error) {
	var compressed bytes.Buffer

	// Encoding the struct using gob
	encoder := gob.NewEncoder(&compressed)
	if err := encoder.Encode(emp); err != nil {
		return nil, fmt.Errorf("failed to encode struct: %w", err)
	}

	// Compressing the encoded data
	w := zlib.NewWriter(&compressed)
	if _, err := w.Write(compressed.Bytes()); err != nil {
		return nil, fmt.Errorf("failed to write compressed data: %w", err)
	}

	if err := w.Close(); err != nil {
		return nil, fmt.Errorf("failed to close zlib writer: %w", err)
	}

	return compressed.Bytes(), nil
}

func decompress(compressedData []byte) (*Emp, error) {
	r, err := zlib.NewReader(bytes.NewReader(compressedData))
	if err != nil {
		return nil, fmt.Errorf("failed to create zlib reader: %w", err)
	}
	defer r.Close()

	// Decompressing the data
	decompressed, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read decompressed data: %w", err)
	}

	// Decoding the struct
	var emp Emp
	decoder := gob.NewDecoder(bytes.NewReader(decompressed))
	if err := decoder.Decode(&emp); err != nil {
		return nil, fmt.Errorf("failed to decode struct: %w", err)
	}

	return &emp, nil
}

func redisHSet(redisClient *redis.Client, key, field string, emp *Emp) error {
	// Compress data before storing
	compressedData, err := compress(emp)
	if err != nil {
		return fmt.Errorf("compression error: %w", err)
	}

	// Store the compressed data in Redis hash field
	if err := redisClient.HSet(context.Background(), key, field, compressedData).Err(); err != nil {
		return fmt.Errorf("error storing data in Redis hash field: %w", err)
	}

	return nil
}

func redisHGet(redisClient *redis.Client, key, field string) (*Emp, error) {
	// Retrieve compressed data from Redis hash field
	retrievedData, err := redisClient.HGet(context.Background(), key, field).Bytes()
	if err != nil {
		return nil, fmt.Errorf("error retrieving data from Redis hash field: %w", err)
	}

	// Decompress and decode the struct
	decompressedEmp, err := decompress(retrievedData)
	if err != nil {
		return nil, fmt.Errorf("decompression error: %w", err)
	}

	return decompressedEmp, nil
}

func main() {
	// Create a Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Your Redis server address
		Password: "",                // No password for local development
		DB:       0,                 // Default DB
	})

	// Example Emp struct with additional fields
	emp := Emp{
		Id:         1,
		IsExist:    true,
		Name:       "John Doe",
		CreatedAt:  time.Now(),
		FloatValue: 123.45,
	}

	// Redis HSet (Hash set)
	hashKey := "your_hash_key"
	field := "your_field"
	if err := redisHSet(redisClient, hashKey, field, &emp); err != nil {
		fmt.Printf("Error inserting data into Redis hash field: %v\n", err)
		return
	}

	// Redis HGet (Hash get)
	retrievedEmp, err := redisHGet(redisClient, hashKey, field)
	if err != nil {
		fmt.Printf("Error reading data from Redis hash field: %v\n", err)
		return
	}

	// Use the retrievedEmp as needed
	fmt.Println("Retrieved Emp struct:", retrievedEmp)
}
