package twitter

import (
	"context"
	"os"
	"testing"
)

// TestIntegration_UploadMediaAndPostTweet tests the complete flow of uploading media and posting a tweet
// This test is skipped by default. To run it, set the environment variable RUN_TWITTER_INTEGRATION_TESTS=true
// and provide valid Twitter API credentials in environment variables.
func TestIntegration_UploadMediaAndPostTweet(t *testing.T) {
	if os.Getenv("RUN_TWITTER_INTEGRATION_TESTS") != "true" {
		t.Skip("Skipping Twitter integration test. Set RUN_TWITTER_INTEGRATION_TESTS=true to run.")
	}

	// Get credentials from environment variables
	apiKey := os.Getenv("TWITTER_API_KEY")
	apiSecret := os.Getenv("TWITTER_API_SECRET")
	accessToken := os.Getenv("TWITTER_ACCESS_TOKEN")
	accessTokenSecret := os.Getenv("TWITTER_ACCESS_TOKEN_SECRET")

	if apiKey == "" || apiSecret == "" || accessToken == "" || accessTokenSecret == "" {
		t.Skip("Skipping integration test: Twitter API credentials not provided")
	}

	// Create client
	client, err := NewClient(apiKey, apiSecret, accessToken, accessTokenSecret)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	// Create a simple test image (1x1 PNG)
	// PNG header + IDAT + IEND
	testImageData := []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, // PNG signature
		0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52, // IHDR chunk
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, // 1x1 image
		0x08, 0x02, 0x00, 0x00, 0x00, 0x90, 0x77, 0x53, // IHDR data
		0xDE, 0x00, 0x00, 0x00, 0x0C, 0x49, 0x44, 0x41, // IDAT chunk
		0x54, 0x08, 0xD7, 0x63, 0xF8, 0xCF, 0xC0, 0x00, // IDAT data
		0x00, 0x03, 0x01, 0x01, 0x00, 0x18, 0xDD, 0x8D, // IDAT data
		0xB4, 0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4E, // IEND chunk
		0x44, 0xAE, 0x42, 0x60, 0x82, // IEND data
	}

	t.Run("upload media", func(t *testing.T) {
		mediaID, err := client.uploadMedia(ctx, testImageData)
		if err != nil {
			t.Fatalf("uploadMedia() failed: %v", err)
		}

		if mediaID == "" {
			t.Error("uploadMedia() returned empty mediaID")
		}

		t.Logf("Successfully uploaded media with ID: %s", mediaID)
	})

	t.Run("post tweet with image", func(t *testing.T) {
		testText := "Test tweet with image from integration test (炎上シミュレーター) #TestOnly"

		result, err := client.PostTweetWithImage(ctx, testText, testImageData)
		if err != nil {
			t.Fatalf("PostTweetWithImage() failed: %v", err)
		}

		if result == nil {
			t.Fatal("PostTweetWithImage() returned nil result")
		}

		if result.ID == "" {
			t.Error("PostTweetWithImage() returned empty ID")
		}

		if result.URL == "" {
			t.Error("PostTweetWithImage() returned empty URL")
		}

		t.Logf("Successfully posted tweet with image: %s", result.URL)
		t.Logf("Tweet ID: %s", result.ID)
	})
}

// TestIntegration_UploadMedia tests only the media upload functionality
func TestIntegration_UploadMedia(t *testing.T) {
	if os.Getenv("RUN_TWITTER_INTEGRATION_TESTS") != "true" {
		t.Skip("Skipping Twitter integration test. Set RUN_TWITTER_INTEGRATION_TESTS=true to run.")
	}

	apiKey := os.Getenv("TWITTER_API_KEY")
	apiSecret := os.Getenv("TWITTER_API_SECRET")
	accessToken := os.Getenv("TWITTER_ACCESS_TOKEN")
	accessTokenSecret := os.Getenv("TWITTER_ACCESS_TOKEN_SECRET")

	if apiKey == "" || apiSecret == "" || accessToken == "" || accessTokenSecret == "" {
		t.Skip("Skipping integration test: Twitter API credentials not provided")
	}

	client, err := NewClient(apiKey, apiSecret, accessToken, accessTokenSecret)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	// Create a simple test image (1x1 PNG)
	testImageData := []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A,
		0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52,
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
		0x08, 0x02, 0x00, 0x00, 0x00, 0x90, 0x77, 0x53,
		0xDE, 0x00, 0x00, 0x00, 0x0C, 0x49, 0x44, 0x41,
		0x54, 0x08, 0xD7, 0x63, 0xF8, 0xCF, 0xC0, 0x00,
		0x00, 0x03, 0x01, 0x01, 0x00, 0x18, 0xDD, 0x8D,
		0xB4, 0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4E,
		0x44, 0xAE, 0x42, 0x60, 0x82,
	}

	mediaID, err := client.uploadMedia(ctx, testImageData)
	if err != nil {
		t.Fatalf("uploadMedia() failed: %v", err)
	}

	if mediaID == "" {
		t.Error("uploadMedia() returned empty mediaID")
	}

	t.Logf("Successfully uploaded media with ID: %s", mediaID)
}
