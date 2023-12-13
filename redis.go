package redis

import (
	"context"
	"digi-model-engine/utils/constants"
	"digi-model-engine/utils/exceptions"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

var (
	RedisClient *redis.Client
)

func InitRedisDB() error {
	// Replace these with your Redis server's connection details.
	redisAddr := viper.GetString("REDIS_ADDR")     // Redis server address
	redisPassword := viper.GetString("REDIS_PASS") // Password (leave empty if no password)
	redisDB := constants.REDIS_DB_NUMBER           // Redis database number

	// Create a Redis client.
	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
	})

	// Ping the Redis server to test the connection.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := client.Ping(ctx).Result(); err != nil {
		return fmt.Errorf("Failed to connect to Redis: %v", err)
	}

	// Assign the Redis client to the global variable for later use.
	RedisClient = client

	return nil
}

func FetchFieldFromRedis(key, field string) string {
	if RedisClient == nil {
		log.Println("Redis client is not initialized.")
		return ""
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	val, err := RedisClient.HGet(ctx, key, field).Result()
	if err != nil {
		if err == redis.Nil {
			exceptions.InternalServerError(err)
		} else {
			log.Fatalf("Error fetching field '%s' from hash '%s' in Redis: %v\n", field, key, err)
		}
		return ""
	}

	return val
}

func InsertFieldIntoRedis(key, field, value string) error {
	if RedisClient == nil {
		return fmt.Errorf("Redis client is not initialized.")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := RedisClient.HSet(ctx, key, field, value).Result()
	if err != nil {
		log.Fatalf("Error inserting field '%s' with value '%s' into hash '%s' in Redis: %v\n", field, value, key, err)
		return err
	}

	return nil
}
