package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
)

func main() {
	// Create Gorilla Mux router
	r := mux.NewRouter()

	// Define route for the HTML page
	r.HandleFunc("/getCount", getCount)
	r.HandleFunc("/", handleHomePage)

	// Start HTTP server
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8081", nil))
}

// Function to handle the home page
func handleHomePage(w http.ResponseWriter, r *http.Request) {
	// Connect to Redis
	ctx := context.Background()
	client := connectToRedis(ctx)

	// Retrieve the number of elements in the list
	redisList := getEnvOrDefault("REDIS_LIST_NAME", "meetup")
	// Count the number of elements in the list
	count64, err := client.LLen(ctx, redisList).Result()
	if err != nil {
		http.Error(w, "Error retrieving number of elements in the list", http.StatusInternalServerError)
		return
	}
	count := int(count64) // Convert int64 to int

	// Generate HTML page
	tmpl := template.Must(template.ParseFiles("index.html"))
	data := struct {
		Count int
	}{
		Count: count,
	}
	tmpl.Execute(w, data)
}

// Function to retrieve the number of elements in the list
func getCount(w http.ResponseWriter, r *http.Request) {
	// Connect to Redis
	ctx := context.Background()
	client := connectToRedis(ctx)

	// Count the number of elements in the list
	redisList := getEnvOrDefault("REDIS_LIST_NAME", "meetup")
	count64, err := client.LLen(ctx, redisList).Result()
	if err != nil {
		http.Error(w, "Error retrieving number of elements in the list", http.StatusInternalServerError)
		return
	}
	count := int(count64) // Convert int64 to int

	// Format the response as JSON
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"count":` + fmt.Sprint(count) + `}`))
}

// Function to connect to Redis
func connectToRedis(ctx context.Context) *redis.Client {
	// Retrieve environment variables
	redisHost := getEnvOrDefault("REDIS_HOST", "localhost")
	redisPort := getEnvOrDefault("REDIS_PORT", "6379")
	redisPassword := os.Getenv("REDIS_PASSWORD")

	// Configure Redis client
	client := redis.NewClient(&redis.Options{
		Addr:     redisHost + ":" + redisPort, // Redis server address
		Password: redisPassword,               // Password, if any
		DB:       0,                           // Database index to select
	})
	return client
}

// Function to get environment variable with default value
func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
