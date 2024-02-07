package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redisbloom/v2"
)

var (
	rdb *redis.Client
	rb  *redisbloom.Client
)

func main() {
	// Connect to Redis
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis server address
		Password: "",               // No password set
		DB:       0,                // Use default DB
	})

	// Connect to RedisBloom
	rb = redisbloom.NewClientFromRedis(rdb)

	// Initialize Gin router
	router := gin.Default()

	// Define routes
	router.GET("/item/:id", getItem)
	router.POST("/item", createItem)
	router.PUT("/item/:id", updateItem)
	router.DELETE("/item/:id", deleteItem)

	// Start server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// getItem retrieves an item from RedisBloom based on ID
func getItem(c *gin.Context) {
	id := c.Param("id")
	exists, err := rb.Exists(id).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if exists == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Item not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

// createItem creates a new item in RedisBloom
func createItem(c *gin.Context) {
	id := c.PostForm("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}
	ok, err := rb.Add(id).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !ok {
		c.JSON(http.StatusConflict, gin.H{"error": "Item already exists"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// updateItem updates an existing item in RedisBloom
func updateItem(c *gin.Context) {
	id := c.Param("id")
	exists, err := rb.Exists(id).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if exists == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Item not found"})
		return
	}
	// Perform update operation here
	c.JSON(http.StatusOK, gin.H{"id": id})
}

// deleteItem deletes an item from RedisBloom
func deleteItem(c *gin.Context) {
	id := c.Param("id")
	deleted, err := rb.Delete(id).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if deleted == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Item not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Item with ID %s deleted", id)})
}
