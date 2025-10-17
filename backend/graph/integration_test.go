package graph

import (
	"context"
	"encoding/base64"
	"os"
	"strings"
	"testing"

	"github.com/Tattsum/enjo/backend/gemini"
	"github.com/Tattsum/enjo/backend/graph/model"
	"github.com/Tattsum/enjo/backend/image"
)

// helper function to check if integration tests should run
func shouldRunIntegrationTests(t *testing.T) (projectID string, location string) {
	t.Helper()
	if os.Getenv("RUN_INTEGRATION_TESTS") != "true" {
		t.Skip("Skipping integration test - set RUN_INTEGRATION_TESTS=true to run")
	}

	projectID = os.Getenv("GCP_PROJECT_ID")
	location = os.Getenv("GCP_LOCATION")

	if projectID == "" {
		t.Skip("Skipping integration test - GCP_PROJECT_ID not set")
	}

	if location == "" {
		location = "us-central1"
	}

	return projectID, location
}

// TestGraphQLImageGenerationIntegration tests the complete GraphQL image generation flow
//
//nolint:revive // Integration test complexity is acceptable
func TestGraphQLImageGenerationIntegration(t *testing.T) {
	projectID, location := shouldRunIntegrationTests(t)
	ctx := context.Background()

	t.Run("complete GraphQL generateImage flow", func(t *testing.T) {
		// Step 1: Initialize clients
		geminiClient, err := gemini.NewClient(ctx, projectID, location)
		if err != nil {
			t.Fatalf("failed to create gemini client: %v", err)
		}
		defer geminiClient.Close()

		imageClient, err := image.NewClient(ctx, projectID, location)
		if err != nil {
			t.Fatalf("failed to create image client: %v", err)
		}
		defer imageClient.Close()

		// Step 2: Create resolver with adapter
		imageAdapter := image.NewAdapter(imageClient)
		resolver := &Resolver{
			geminiClient: geminiClient,
			imageClient:  imageAdapter,
		}

		// Step 3: Prepare input
		text := "今日のランチはラーメンでした！美味しかった！"
		input := model.GenerateImageInput{
			Text: text,
		}

		// Step 4: Execute generateImage mutation
		result, err := resolver.Mutation().GenerateImage(ctx, input)
		if err != nil {
			t.Fatalf("generateImage mutation failed: %v", err)
		}

		// Step 5: Validate result
		if result == nil {
			t.Fatal("expected result to be non-nil")
		}

		if result.ImageURL == "" {
			t.Error("expected imageUrl to be non-empty")
		}

		// Validate data URL format
		if !strings.HasPrefix(result.ImageURL, "data:image/png;base64,") {
			t.Errorf("expected data URL to start with 'data:image/png;base64,', got: %s", result.ImageURL[:50])
		}

		// Decode and validate image data
		base64Data := strings.TrimPrefix(result.ImageURL, "data:image/png;base64,")
		imageData, err := base64.StdEncoding.DecodeString(base64Data)
		if err != nil {
			t.Fatalf("failed to decode base64 image data: %v", err)
		}

		if len(imageData) == 0 {
			t.Error("expected decoded image data to be non-empty")
		}

		// Verify PNG header
		pngHeader := []byte{0x89, 0x50, 0x4E, 0x47}
		if len(imageData) >= 4 {
			for i, b := range pngHeader {
				if imageData[i] != b {
					t.Errorf("invalid PNG header at byte %d: expected %x, got %x", i, b, imageData[i])
				}
			}
		}

		if result.Prompt == "" {
			t.Error("expected prompt to be non-empty")
		}

		if result.GeneratedAt == "" {
			t.Error("expected generatedAt to be non-empty")
		}

		t.Log("Successfully generated image via GraphQL")
		t.Logf("- Image size: %d bytes", len(imageData))
		t.Logf("- Prompt: %s", result.Prompt)
		t.Logf("- Generated at: %s", result.GeneratedAt)
	})

	t.Run("generateImage with different styles", func(t *testing.T) {
		geminiClient, err := gemini.NewClient(ctx, projectID, location)
		if err != nil {
			t.Fatalf("failed to create gemini client: %v", err)
		}
		defer geminiClient.Close()

		imageClient, err := image.NewClient(ctx, projectID, location)
		if err != nil {
			t.Fatalf("failed to create image client: %v", err)
		}
		defer imageClient.Close()

		imageAdapter := image.NewAdapter(imageClient)
		resolver := &Resolver{
			geminiClient: geminiClient,
			imageClient:  imageAdapter,
		}

		testCases := []struct {
			name  string
			style *model.ImageStyle
		}{
			{
				name:  "default style (nil)",
				style: nil,
			},
			{
				name: "MEME style",
				style: func() *model.ImageStyle {
					s := model.ImageStyleMeme
					return &s
				}(),
			},
			{
				name: "REALISTIC style",
				style: func() *model.ImageStyle {
					s := model.ImageStyleRealistic
					return &s
				}(),
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				input := model.GenerateImageInput{
					Text:  "シンプルなテスト",
					Style: tc.style,
				}

				result, err := resolver.Mutation().GenerateImage(ctx, input)
				if err != nil {
					t.Fatalf("generateImage failed: %v", err)
				}

				if result.ImageURL == "" {
					t.Error("expected imageUrl to be non-empty")
				}

				// Extract and validate image size
				base64Data := strings.TrimPrefix(result.ImageURL, "data:image/png;base64,")
				imageData, err := base64.StdEncoding.DecodeString(base64Data)
				if err != nil {
					t.Fatalf("failed to decode base64: %v", err)
				}

				t.Logf("Style %s: generated %d bytes", tc.name, len(imageData))
			})
		}
	})

	t.Run("generateImage with different aspect ratios", func(t *testing.T) {
		geminiClient, err := gemini.NewClient(ctx, projectID, location)
		if err != nil {
			t.Fatalf("failed to create gemini client: %v", err)
		}
		defer geminiClient.Close()

		imageClient, err := image.NewClient(ctx, projectID, location)
		if err != nil {
			t.Fatalf("failed to create image client: %v", err)
		}
		defer imageClient.Close()

		imageAdapter := image.NewAdapter(imageClient)
		resolver := &Resolver{
			geminiClient: geminiClient,
			imageClient:  imageAdapter,
		}

		testCases := []struct {
			name        string
			aspectRatio *model.AspectRatio
		}{
			{
				name:        "default (nil)",
				aspectRatio: nil,
			},
			{
				name: "SQUARE",
				aspectRatio: func() *model.AspectRatio {
					ar := model.AspectRatioSquare
					return &ar
				}(),
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				input := model.GenerateImageInput{
					Text:        "テスト画像",
					AspectRatio: tc.aspectRatio,
				}

				result, err := resolver.Mutation().GenerateImage(ctx, input)
				if err != nil {
					t.Fatalf("generateImage failed: %v", err)
				}

				if result.ImageURL == "" {
					t.Error("expected imageUrl to be non-empty")
				}

				t.Logf("AspectRatio %s: success", tc.name)
			})
		}
	})

	t.Run("error handling - empty text", func(t *testing.T) {
		geminiClient, err := gemini.NewClient(ctx, projectID, location)
		if err != nil {
			t.Fatalf("failed to create gemini client: %v", err)
		}
		defer geminiClient.Close()

		imageClient, err := image.NewClient(ctx, projectID, location)
		if err != nil {
			t.Fatalf("failed to create image client: %v", err)
		}
		defer imageClient.Close()

		imageAdapter := image.NewAdapter(imageClient)
		resolver := &Resolver{
			geminiClient: geminiClient,
			imageClient:  imageAdapter,
		}

		input := model.GenerateImageInput{
			Text: "",
		}

		_, err = resolver.Mutation().GenerateImage(ctx, input)
		if err == nil {
			t.Error("expected error for empty text input")
		}

		expectedMsg := "text is required"
		if !strings.Contains(err.Error(), expectedMsg) {
			t.Errorf("expected error to contain %q, got: %v", expectedMsg, err)
		}
	})
}

// TestGraphQLImageGenerationWithTwitterIntegration tests the complete flow
// from image generation to Twitter posting
func TestGraphQLImageGenerationWithTwitterIntegration(t *testing.T) {
	projectID, location := shouldRunIntegrationTests(t)
	ctx := context.Background()

	// This test is more complex as it requires Twitter API credentials
	// For now, we'll test the preparation of the data structure
	t.Run("prepare data for Twitter posting", func(t *testing.T) {
		geminiClient, err := gemini.NewClient(ctx, projectID, location)
		if err != nil {
			t.Fatalf("failed to create gemini client: %v", err)
		}
		defer geminiClient.Close()

		imageClient, err := image.NewClient(ctx, projectID, location)
		if err != nil {
			t.Fatalf("failed to create image client: %v", err)
		}
		defer imageClient.Close()

		imageAdapter := image.NewAdapter(imageClient)
		resolver := &Resolver{
			geminiClient: geminiClient,
			imageClient:  imageAdapter,
		}

		// Step 1: Generate image
		input := model.GenerateImageInput{
			Text: "統合テスト用の投稿",
		}

		result, err := resolver.Mutation().GenerateImage(ctx, input)
		if err != nil {
			t.Fatalf("generateImage failed: %v", err)
		}

		// Step 2: Prepare Twitter post input with image
		twitterInput := model.TwitterPostInput{
			Text:     "統合テスト用の投稿",
			ImageURL: &result.ImageURL,
		}

		// Step 3: Validate the data structure is ready for Twitter posting
		if twitterInput.Text == "" {
			t.Error("expected text to be non-empty")
		}

		if twitterInput.ImageURL == nil || *twitterInput.ImageURL == "" {
			t.Error("expected imageUrl to be non-empty")
		}

		// Verify image URL is a valid data URL
		if !strings.HasPrefix(*twitterInput.ImageURL, "data:image/png;base64,") {
			t.Error("expected valid data URL format")
		}

		t.Log("Successfully prepared data for Twitter posting with image")
		t.Logf("- Text: %s", twitterInput.Text)
		t.Logf("- Image URL length: %d", len(*twitterInput.ImageURL))
	})
}
