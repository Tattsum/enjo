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
)

// setupRouter creates and configures the HTTP router
func setupRouter(geminiClient graph.GeminiClient) http.Handler {
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
	resolver := graph.NewResolver(geminiClient)

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

	// Get Gemini API key
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("GEMINI_API_KEY environment variable is required")
	}

	// Get port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize Gemini client
	ctx := context.Background()
	geminiClient, err := gemini.NewClient(ctx, apiKey)
	if err != nil {
		log.Fatalf("Failed to create Gemini client: %v", err)
	}

	// Setup router
	router := setupRouter(geminiClient)

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
