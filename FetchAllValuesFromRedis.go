package main

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/go-redis/redis/v8"
)

// FetchAllValues fetches all values from a Redis hash and converts them to their proper Go types.
func FetchAllValues(redisClient *redis.Client, hashKey string) (map[string]interface{}, error) {
	ctx := context.Background()

	// Fetch all fields and values from the hash
	fieldsAndValues, err := redisClient.HGetAll(ctx, hashKey).Result()
	if err != nil {
		return nil, err
	}

	// Convert values to proper Go types
	result := make(map[string]interface{})
	for field, value := range fieldsAndValues {
		convertedValue, err := convertToGoType(value)
		if err != nil {
			return nil, err
		}
		result[field] = convertedValue
	}

	return result, nil
}

// convertToGoType converts a string value from Redis to its proper Go type.
func convertToGoType(value string) (interface{}, error) {
	// You can extend this function to handle more types as needed.
	// Currently, it supports converting to int, float64, and string.

	// Try converting to int
	if intValue, err := strconv.Atoi(value); err == nil {
		return intValue, nil
	}

	// Try converting to float64
	if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
		return floatValue, nil
	}

	// If no conversion is successful, return the original string value
	return value, nil
}

func main() {
	// Replace these values with your actual Redis server information
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // No password
		DB:       0,  // Default DB
	})

	// Close the connection when the function exits
	defer func() {
		if err := redisClient.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	// Replace "yourHashKey" with the actual hash key you are using
	hashKey := "yourHashKey"

	// Example usage
	details, err := FetchAllValues(redisClient, hashKey)
	if err != nil {
		log.Fatal(err)
	}

	// Print the fetched details and their types
	for field, value := range details {
		fmt.Printf("%s: %v (Type: %T)\n", field, value, value)
	}
}
