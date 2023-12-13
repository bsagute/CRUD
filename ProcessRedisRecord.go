package main

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/go-redis/redis/v8"
)

// InsertRecordInPostgress is a placeholder function. You should replace it with your actual logic.
func InsertRecordInPostgress(hashKey string, details map[string]interface{}) {
	// Replace this with your actual implementation
	fmt.Printf("Inserting record in Postgres for hashKey: %s with details: %+v\n", hashKey, details)
	// Implement your logic to insert the record into Postgres here
}

func processRecord(redisClient *redis.Client, hashKey string) error {
	ctx := context.Background()

	// Check if the hashKey exists in Redis
	exists, err := redisClient.Exists(ctx, hashKey).Result()
	if err != nil {
		return err
	}

	// Get the IsInserted flag value
	isInserted, err := redisClient.HGet(ctx, hashKey, "IsInserted").Result()
	if err != nil && err != redis.Nil {
		return err
	}

	// Convert isInserted to a boolean
	isInsertedFlag, err := strconv.ParseBool(isInserted)
	if err != nil {
		return err
	}

	if exists == 0 {
		// Condition 1: HashKey doesn't exist, create a new record in Redis (IsInserted set to false)
		recordDetails := map[string]interface{}{
			// Add your details here
			"IsInserted": false, // Set the flag to false when creating a new record
			// Add other details as needed
		}

		// Set the new record in Redis
		err := redisClient.HMSet(ctx, hashKey, recordDetails).Err()
		if err != nil {
			return err
		}

		fmt.Println("Condition 1: New record created in Redis (IsInserted set to false).")
	} else if isInsertedFlag {
		// Condition 2: If IsInserted flag is true, do nothing and return a message
		fmt.Println("Condition 2: Input hash key exists, and IsInserted flag is true. Nothing to do.")
	} else {
		// Condition 3: If IsInserted flag is false, update the flag and call InsertRecordInPostgress
		err := redisClient.HSet(ctx, hashKey, "IsInserted", true).Err()
		if err != nil {
			return err
		}

		// Fetch the details from Redis
		details, err := FetchAllValues(redisClient, hashKey)
		if err != nil {
			return err
		}

		// Call the function to insert record in Postgres
		InsertRecordInPostgress(hashKey, details)

		fmt.Println("Condition 3: IsInserted flag updated to true, and record inserted into Postgres.")
	}

	return nil
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
	err := processRecord(redisClient, hashKey)
	if err != nil {
		log.Fatal(err)
	}
}
