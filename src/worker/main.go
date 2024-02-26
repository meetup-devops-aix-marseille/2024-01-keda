package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

func main() {
	// Retrieving environment variables with default values
	redisHost := getEnvOrDefault("REDIS_HOST", "localhost")
	redisPort := getEnvOrDefault("REDIS_PORT", "6379")
	redisPassword := os.Getenv("REDIS_PASSWORD")

	// Configuring Redis client
	client := redis.NewClient(&redis.Options{
		Addr:     redisHost + ":" + redisPort, // Redis server address
		Password: redisPassword,               // Password, if any
		DB:       0,                           // Database index to select
	})

	// Context for Redis operations
	ctx := context.Background()

	// Redis list key
	redisList := getEnvOrDefault("REDIS_LIST_NAME", "meetup")

	// Check the connection with Redis
	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Error connecting to Redis: %v", err)
	}

	// Count the number of elements in the list
	value, err := client.LLen(ctx, redisList).Result()
	if err != nil {
		log.Printf("Error reading the list: %v", err)
		return
	}

	// Value of the list
	log.Printf("Current value of the list: %d", value)

	// If the current value is 0, do nothing and exit
	if value == 0 {
		log.Println("The list is empty, exiting the application")
		return
	}

	// Wait for 10 seconds
	log.Println("Waiting for 10 seconds")
	time.Sleep(10 * time.Second)

	// Remove an item from the list
	_, err = client.LPop(ctx, redisList).Result()
	if err != nil {
		log.Printf("Error decrementing the list: %v", err)
		return
	}

	// End of the application
	log.Println("End of the application")
}

// Function to get environment variable with default value
func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
