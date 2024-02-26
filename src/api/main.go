package main

import (
	"context"
	"log"
	"net/http"
	"os"

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

	// Add a task to the list
	http.HandleFunc("/add-task", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Not allowed method", http.StatusMethodNotAllowed)
			return
		}

		// Add a new element to the list
		value, err := client.RPush(ctx, redisList, "task").Result()
		if err != nil {
			log.Printf("Error incrementing the list: %v", err)
			http.Error(w, "Error incrementing the list", http.StatusInternalServerError)
			return
		}

		// New value of the list
		log.Printf("New value of the list: %d", value)
		w.WriteHeader(http.StatusCreated)
	})

	// Start the HTTP server
	log.Println("Server started on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Function to get environment variable with default value
func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
