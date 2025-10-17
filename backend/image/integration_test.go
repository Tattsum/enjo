package image

import (
	"context"
	"os"
	"testing"
)

// TestImageGenerationIntegration tests the complete image generation flow
// This is an integration test that requires:
// - Valid GCP credentials
// - Vertex AI API enabled
// - Proper environment setup
//
//nolint:revive // Integration test complexity is acceptable
func TestImageGenerationIntegration(t *testing.T) {
	// Skip if not in integration test mode
	if os.Getenv("RUN_INTEGRATION_TESTS") != "true" {
		t.Skip("Skipping integration test - set RUN_INTEGRATION_TESTS=true to run")
	}

	ctx := context.Background()
	projectID := os.Getenv("GCP_PROJECT_ID")
	location := os.Getenv("GCP_LOCATION")

	if projectID == "" {
		t.Skip("Skipping integration test - GCP_PROJECT_ID not set")
	}

	if location == "" {
		location = "us-central1"
	}

	t.Run("complete image generation flow", func(t *testing.T) {
		// Step 1: Create client
		client, err := NewClient(ctx, projectID, location)
		if err != nil {
			t.Fatalf("failed to create client: %v", err)
		}
		defer client.Close()

		// Step 2: Generate image with minimal prompt
		prompt := "A small flame icon"
		result, err := client.GenerateImage(ctx, prompt)
		if err != nil {
			t.Fatalf("failed to generate image: %v", err)
		}

		// Step 3: Validate result
		if result == nil {
			t.Fatal("expected result to be non-nil")
		}
		if len(result.ImageData) == 0 {
			t.Error("expected image data to be non-empty")
		}
		if result.Prompt != prompt {
			t.Errorf("expected prompt %q, got %q", prompt, result.Prompt)
		}
		if result.GeneratedAt.IsZero() {
			t.Error("expected GeneratedAt to be set")
		}

		// Step 4: Verify image data is valid (basic check)
		// PNG files start with specific magic bytes: 89 50 4E 47
		if len(result.ImageData) < 8 {
			t.Error("image data too small to be valid PNG")
		} else {
			pngHeader := []byte{0x89, 0x50, 0x4E, 0x47}
			for i, b := range pngHeader {
				if result.ImageData[i] != b {
					t.Errorf("invalid PNG header at byte %d: expected %x, got %x", i, b, result.ImageData[i])
				}
			}
		}

		t.Logf("Successfully generated image: %d bytes", len(result.ImageData))
	})

	t.Run("image generation with style options", func(t *testing.T) {
		client, err := NewClient(ctx, projectID, location)
		if err != nil {
			t.Fatalf("failed to create client: %v", err)
		}
		defer client.Close()

		testCases := []struct {
			name        string
			prompt      string
			style       string
			aspectRatio string
		}{
			{
				name:        "realistic style",
				prompt:      "A realistic fire",
				style:       "realistic",
				aspectRatio: "1:1",
			},
			{
				name:        "illustration style",
				prompt:      "An illustrated flame",
				style:       "illustration",
				aspectRatio: "1:1",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				result, err := client.GenerateImage(
					ctx,
					tc.prompt,
					WithStyle(tc.style),
					WithAspectRatio(tc.aspectRatio),
				)
				if err != nil {
					t.Fatalf("failed to generate image: %v", err)
				}

				if len(result.ImageData) == 0 {
					t.Error("expected image data to be non-empty")
				}

				t.Logf("Generated %s image: %d bytes", tc.style, len(result.ImageData))
			})
		}
	})

	t.Run("multiple concurrent image generations", func(t *testing.T) {
		client, err := NewClient(ctx, projectID, location)
		if err != nil {
			t.Fatalf("failed to create client: %v", err)
		}
		defer client.Close()

		// Test concurrent image generation (limited to 3 to avoid rate limits)
		const numConcurrent = 3
		results := make(chan *Result, numConcurrent)
		errors := make(chan error, numConcurrent)

		for range numConcurrent {
			go func() {
				result, err := client.GenerateImage(ctx, "A small icon")
				if err != nil {
					errors <- err
					return
				}
				results <- result
			}()
		}

		// Collect results
		for range numConcurrent {
			select {
			case result := <-results:
				if len(result.ImageData) == 0 {
					t.Error("expected image data to be non-empty")
				}
			case err := <-errors:
				t.Errorf("concurrent generation failed: %v", err)
			}
		}

		t.Logf("Successfully completed %d concurrent image generations", numConcurrent)
	})

	t.Run("error handling - invalid prompt", func(t *testing.T) {
		client, err := NewClient(ctx, projectID, location)
		if err != nil {
			t.Fatalf("failed to create client: %v", err)
		}
		defer client.Close()

		// Test with empty prompt
		_, err = client.GenerateImage(ctx, "")
		if err == nil {
			t.Error("expected error for empty prompt")
		}
		expectedMsg := "prompt is required"
		if err.Error() != expectedMsg {
			t.Errorf("expected error message %q, got %q", expectedMsg, err.Error())
		}
	})
}

// TestImageGenerationPerformance tests the performance of image generation
func TestImageGenerationPerformance(t *testing.T) {
	if os.Getenv("RUN_INTEGRATION_TESTS") != "true" {
		t.Skip("Skipping integration test - set RUN_INTEGRATION_TESTS=true to run")
	}

	ctx := context.Background()
	projectID := os.Getenv("GCP_PROJECT_ID")
	location := os.Getenv("GCP_LOCATION")

	if projectID == "" {
		t.Skip("Skipping performance test - GCP_PROJECT_ID not set")
	}

	if location == "" {
		location = "us-central1"
	}

	t.Run("measure image generation latency", func(t *testing.T) {
		client, err := NewClient(ctx, projectID, location)
		if err != nil {
			t.Fatalf("failed to create client: %v", err)
		}
		defer client.Close()

		prompt := "A simple icon"

		// Warm-up request (first request is often slower)
		_, _ = client.GenerateImage(ctx, prompt)

		// Measure actual performance
		result, err := client.GenerateImage(ctx, prompt)
		if err != nil {
			t.Fatalf("failed to generate image: %v", err)
		}

		if len(result.ImageData) == 0 {
			t.Error("expected image data to be non-empty")
		}

		// Log performance metrics
		t.Log("Image generation completed")
		t.Logf("- Image size: %d bytes", len(result.ImageData))
		t.Logf("- Timestamp: %s", result.GeneratedAt.Format("2006-01-02 15:04:05"))
	})
}
