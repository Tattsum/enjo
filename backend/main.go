package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	"github.com/Tattsum/enjo/backend/gemini"
	"github.com/Tattsum/enjo/backend/graph"
	"github.com/Tattsum/enjo/backend/graph/generated"
	"github.com/Tattsum/enjo/backend/twitter"
)

// setupRouter creates and configures the HTTP router
func setupRouter(geminiClient graph.GeminiClient, twitterClient graph.TwitterClient) http.Handler {
	router := chi.NewRouter()

	// CORS configuration
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:8080"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Health check endpoint
	router.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(map[string]string{"status": "OK"}); err != nil {
			log.Printf("Error encoding health response: %v", err)
		}
	})

	// GraphQL resolver
	resolver := graph.NewResolver(geminiClient, twitterClient)

	// GraphQL server
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
		Resolvers: resolver,
	}))

	// GraphQL endpoints
	router.Handle("/graphql", srv)
	router.Handle("/", playground.Handler("GraphQL Playground", "/graphql"))

	return router
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	// Get GCP configuration
	projectID := os.Getenv("GCP_PROJECT_ID")
	if projectID == "" {
		log.Fatal("GCP_PROJECT_ID environment variable is required")
	}

	location := os.Getenv("GCP_LOCATION")
	if location == "" {
		location = "us-central1" // Default location
	}

	// Get port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize Vertex AI client with ADC
	ctx := context.Background()
	geminiClient, err := gemini.NewClient(ctx, projectID, location)
	if err != nil {
		log.Fatalf("Failed to create Vertex AI client: %v", err)
	}

	// Initialize Twitter client
	twitterAPIKey := os.Getenv("TWITTER_API_KEY")
	twitterAPISecret := os.Getenv("TWITTER_API_SECRET")
	twitterAccessToken := os.Getenv("TWITTER_ACCESS_TOKEN")
	twitterAccessTokenSecret := os.Getenv("TWITTER_ACCESS_TOKEN_SECRET")

	// Twitter client is optional - if not configured, operations will fail gracefully
	var twitterClient graph.TwitterClient
	if twitterAPIKey != "" && twitterAPISecret != "" && twitterAccessToken != "" && twitterAccessTokenSecret != "" {
		twitterClient, err = twitter.NewClient(twitterAPIKey, twitterAPISecret, twitterAccessToken, twitterAccessTokenSecret)
		if err != nil {
			log.Printf("Warning: Failed to create Twitter client: %v", err)
			log.Println("Twitter posting functionality will be disabled")
		} else {
			log.Println("Twitter client initialized successfully")
		}
	} else {
		log.Println("Twitter API credentials not configured - Twitter posting functionality will be disabled")
	}

	// Setup router
	router := setupRouter(geminiClient, twitterClient)

	// Start server
	log.Printf("Server is running on http://localhost:%s", port)
	log.Printf("GraphQL Playground: http://localhost:%s/graphql", port)

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}
