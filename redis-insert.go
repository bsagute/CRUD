package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/spf13/viper"

	"your-package-path/redis" // Replace with the actual package path
)

func main() {
	// Initialize Redis client
	err := redis.InitRedisDB()
	if err != nil {
		log.Fatalf("Failed to initialize Redis client: %v", err)
	}

	// JSON data to insert into Redis
	jsonData := `
		[
			{
				"metrics_status": "TODO",
				"entity_id": "TODO",
				"redis_key": "298815d6-dc3e-4f2a-a4b3-e8934c88bdbf",
				"created_at": "2023-12-06 16:51:58"
			},
			{
				"metrics_status": "TODO",
				"entity_id": "TODO",
				"redis_key": "0f71fa3f-1ffb-411d-ae61-ac8073a1f13f",
				"created_at": "2023-12-07 11:59:02"
			},
			{
				"metrics_status": "TODO",
				"entity_id": "TODO",
				"redis_key": "4813b36e-dba5-47a1-aaad-f1a6d04e5957",
				"created_at": "2023-12-07 13:27:44"
			}
		]
	`

	// Parse JSON data into a slice of structs
	var data []struct {
		MetricsStatus string `json:"metrics_status"`
		EntityID      string `json:"entity_id"`
		RedisKey      string `json:"redis_key"`
		CreatedAt     string `json:"created_at"`
	}
	if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
		log.Fatalf("Failed to parse JSON data: %v", err)
	}

	// Insert data into Redis
	for _, item := range data {
		value := map[string]interface{}{
			"metrics_status": item.MetricsStatus,
			"entity_id":      item.EntityID,
			"created_at":     item.CreatedAt,
		}

		// Convert value to JSON before inserting (you can customize this based on your needs)
		jsonValue, err := json.Marshal(value)
		if err != nil {
			log.Printf("Failed to convert value to JSON: %v", err)
			continue
		}

		// Insert key-value pair into Redis
		err = redis.InsertKeyToRedis(item.RedisKey, string(jsonValue))
		if err != nil {
			log.Printf("Failed to insert data for key '%s' into Redis: %v", item.RedisKey, err)
		} else {
			log.Printf("Data for key '%s' inserted successfully into Redis", item.RedisKey)
		}
	}
}
