package redis

import (
	"context"
	"digi-model-engine/utils/constants"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

var (
	RedisClient *redis.Client
)
// Load Environment Variables
	env := os.Getenv("APP_ENV")
	envPath := "configs/" + env + ".env"

	viper.SetConfigFile(envPath)
	viper.ReadInConfig()
	viper.AutomaticEnv()
	redis.InitRedisDB()
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

func FetchKeyFromRedis(key string) string {
	if RedisClient == nil {
		log.Println("Redis client is not initialized.")
		return ""
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	val, err := RedisClient.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			fmt.Printf("Key '%s' does not exist in Redis\n", key)
		} else {
			log.Fatalf("Error fetching key '%s' from Redis: %v\n", key, err)
		}
		return ""
	}

	fmt.Printf("Key '%s' has the value: %s\n", key, val)
	return val
}
