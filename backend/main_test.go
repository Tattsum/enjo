package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Tattsum/enjo/backend/graph"
)

// MockGeminiClient for testing
type MockGeminiClient struct{}

func (*MockGeminiClient) GenerateInflammatoryText(_ context.Context, _ string, _ int) (string, error) {
	return "Mock inflammatory text", nil
}

func (*MockGeminiClient) GenerateExplanation(_ context.Context, _, _ string) (string, error) {
	return "Mock explanation", nil
}

func (*MockGeminiClient) GenerateReply(_ context.Context, _, _ string) (string, error) {
	return "Mock reply", nil
}

func TestHealthEndpoint(t *testing.T) {
	// Arrange
	handler := setupRouter(&MockGeminiClient{})
	req := httptest.NewRequest(http.MethodGet, "/health", http.NoBody)
	w := httptest.NewRecorder()

	// Act
	handler.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	expected := `{"status":"OK"}`
	actual := strings.TrimSpace(w.Body.String())
	if actual != expected {
		t.Errorf("Expected body %q, got %q", expected, actual)
	}
}

func TestGraphQLEndpoint(t *testing.T) {
	// Arrange
	handler := setupRouter(&MockGeminiClient{})

	// GraphQL health query
	query := `{"query": "query { health }"}`
	req := httptest.NewRequest(http.MethodPost, "/graphql", strings.NewReader(query))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	handler.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	data, ok := response["data"].(map[string]any)
	if !ok {
		t.Fatal("Expected 'data' field in response")
	}

	health, ok := data["health"].(string)
	if !ok {
		t.Fatal("Expected 'health' field in data")
	}

	if health != "OK" {
		t.Errorf("Expected health to be 'OK', got %q", health)
	}
}

func TestCORSHeaders(t *testing.T) {
	// Arrange
	handler := setupRouter(&MockGeminiClient{})
	req := httptest.NewRequest(http.MethodOptions, "/graphql", http.NoBody)
	req.Header.Set("Origin", "http://localhost:3000")
	req.Header.Set("Access-Control-Request-Method", "POST")
	w := httptest.NewRecorder()

	// Act
	handler.ServeHTTP(w, req)

	// Assert
	allowOrigin := w.Header().Get("Access-Control-Allow-Origin")
	if allowOrigin == "" {
		t.Error("Expected Access-Control-Allow-Origin header to be set")
	}

	allowMethods := w.Header().Get("Access-Control-Allow-Methods")
	if allowMethods == "" {
		t.Error("Expected Access-Control-Allow-Methods header to be set")
	}
}

func TestSetupRouter(t *testing.T) {
	// Arrange
	client := &MockGeminiClient{}

	// Act
	handler := setupRouter(client)

	// Assert
	if handler == nil {
		t.Fatal("Expected setupRouter to return non-nil handler")
	}
}

func TestNewResolver(t *testing.T) {
	// Arrange
	client := &MockGeminiClient{}

	// Act
	resolver := graph.NewResolver(client)

	// Assert
	if resolver == nil {
		t.Fatal("Expected NewResolver to return non-nil resolver")
	}
}
