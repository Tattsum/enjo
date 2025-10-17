package image

import (
	"context"
	"testing"
)

func TestNewClient(t *testing.T) {
	ctx := context.Background()

	t.Run("error when projectID is empty", func(t *testing.T) {
		_, err := NewClient(ctx, "", "us-central1")
		if err == nil {
			t.Fatal("expected error when projectID is empty")
		}
		expectedMsg := "GCP project ID is required"
		if err.Error() != expectedMsg {
			t.Errorf("expected error message %q, got %q", expectedMsg, err.Error())
		}
	})

	t.Run("successful client creation - integration test", func(t *testing.T) {
		// Skip this test in regular unit test runs
		// This is an integration test that requires real Vertex AI credentials
		t.Skip("Skipping integration test - requires real Vertex AI credentials and API access")

		projectID := "test-project"
		location := "us-central1"

		client, err := NewClient(ctx, projectID, location)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if client == nil {
			t.Fatal("expected client to be non-nil")
		}

		defer client.Close()

		if client.projectID != projectID {
			t.Errorf("expected projectID %s, got %s", projectID, client.projectID)
		}
		if client.location != location {
			t.Errorf("expected location %s, got %s", location, client.location)
		}
	})

	t.Run("uses default location when empty - integration test", func(t *testing.T) {
		// Skip this test in regular unit test runs
		t.Skip("Skipping integration test - requires real Vertex AI credentials and API access")

		projectID := "test-project"
		location := ""

		client, err := NewClient(ctx, projectID, location)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if client == nil {
			t.Fatal("expected client to be non-nil")
		}

		defer client.Close()

		expectedLocation := "us-central1"
		if client.location != expectedLocation {
			t.Errorf("expected location %s, got %s", expectedLocation, client.location)
		}
	})
}

func TestGenerateImage(t *testing.T) {
	ctx := context.Background()

	t.Run("error when prompt is empty", func(t *testing.T) {
		client := createTestClient(t, ctx)
		defer client.Close()

		_, err := client.GenerateImage(ctx, "")
		if err == nil {
			t.Fatal("expected error when prompt is empty")
		}
		expectedMsg := "prompt is required"
		if err.Error() != expectedMsg {
			t.Errorf("expected error message %q, got %q", expectedMsg, err.Error())
		}
	})

	t.Run("successful image generation - integration test", func(t *testing.T) {
		// Skip this test in regular unit test runs
		// This is an integration test that requires real Vertex AI credentials
		t.Skip("Skipping integration test - requires real Vertex AI credentials and API access")

		client := createTestClient(t, ctx)
		defer client.Close()

		prompt := "A dramatic image of fire and controversy on social media"
		result := generateTestImage(t, ctx, client, prompt)
		validateImageResult(t, result, prompt)
	})

	t.Run("image generation with options - integration test", func(t *testing.T) {
		// Skip this test in regular unit test runs
		t.Skip("Skipping integration test - requires real Vertex AI credentials and API access")

		client := createTestClient(t, ctx)
		defer client.Close()

		prompt := "A dramatic image of fire"
		result := generateTestImageWithOptions(t, ctx, client, prompt)
		if len(result.ImageData) == 0 {
			t.Error("expected image data to be non-empty")
		}
	})
}

func createTestClient(t *testing.T, ctx context.Context) *Client {
	t.Helper()
	client, err := NewClient(ctx, "test-project", "us-central1")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	return client
}

func generateTestImage(t *testing.T, ctx context.Context, client *Client, prompt string) *Result {
	t.Helper()
	result, err := client.GenerateImage(ctx, prompt)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("expected result to be non-nil")
	}
	return result
}

func generateTestImageWithOptions(t *testing.T, ctx context.Context, client *Client, prompt string) *Result {
	t.Helper()
	result, err := client.GenerateImage(ctx, prompt, WithAspectRatio("1:1"), WithStyle("realistic"))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("expected result to be non-nil")
	}
	return result
}

func validateImageResult(t *testing.T, result *Result, expectedPrompt string) {
	t.Helper()
	if len(result.ImageData) == 0 {
		t.Error("expected image data to be non-empty")
	}
	if result.Prompt != expectedPrompt {
		t.Errorf("expected prompt %s, got %s", expectedPrompt, result.Prompt)
	}
	if result.GeneratedAt.IsZero() {
		t.Error("expected GeneratedAt to be set")
	}
}

func TestImageOptions(t *testing.T) {
	t.Run("WithStyle option", func(t *testing.T) {
		opts := &imageOptions{}
		WithStyle("realistic")(opts)
		if opts.style != "realistic" {
			t.Errorf("expected style 'realistic', got %s", opts.style)
		}
	})

	t.Run("WithAspectRatio option", func(t *testing.T) {
		opts := &imageOptions{}
		WithAspectRatio("16:9")(opts)
		if opts.aspectRatio != "16:9" {
			t.Errorf("expected aspectRatio '16:9', got %s", opts.aspectRatio)
		}
	})

	t.Run("WithSize option", func(t *testing.T) {
		opts := &imageOptions{}
		WithSize(1024, 1024)(opts)
		if opts.width != 1024 {
			t.Errorf("expected width 1024, got %d", opts.width)
		}
		if opts.height != 1024 {
			t.Errorf("expected height 1024, got %d", opts.height)
		}
	})
}
