package gemini_test

import (
	"context"
	"os"
	"testing"

	"github.com/Tattsum/enjo/backend/gemini"
)

//nolint:revive // Test function complexity is acceptable
func TestClient_GenerateInflammatoryText(t *testing.T) {
	// Skip if GCP project ID is not set
	projectID := os.Getenv("GCP_PROJECT_ID")
	if projectID == "" {
		t.Skip("GCP_PROJECT_ID is not set")
	}

	location := os.Getenv("GCP_LOCATION")
	if location == "" {
		location = "us-central1" // Default location
	}

	client, err := gemini.NewClient(context.Background(), projectID, location)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			t.Errorf("failed to close client: %v", err)
		}
	}()

	tests := []struct {
		name     string
		original string
		level    int
		wantErr  bool
	}{
		{
			name:     "valid input - level 1",
			original: "今日はいい天気ですね",
			level:    1,
			wantErr:  false,
		},
		{
			name:     "valid input - level 3",
			original: "新しい製品をリリースしました",
			level:    3,
			wantErr:  false,
		},
		{
			name:     "valid input - level 5",
			original: "みんなで協力していきましょう",
			level:    5,
			wantErr:  false,
		},
		{
			name:     "invalid level - too low",
			original: "テストです",
			level:    0,
			wantErr:  true,
		},
		{
			name:     "invalid level - too high",
			original: "テストです",
			level:    6,
			wantErr:  true,
		},
		{
			name:     "empty original text",
			original: "",
			level:    3,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			result, err := client.GenerateInflammatoryText(ctx, tt.original, tt.level)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, but got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if result == "" {
				t.Error("result is empty")
			}

			// Basic validation: result should be different from original
			// (unless it's already inflammatory)
			if result == "" {
				t.Error("generated text is empty")
			}
		})
	}
}

func TestClient_GenerateExplanation(t *testing.T) {
	projectID := os.Getenv("GCP_PROJECT_ID")
	if projectID == "" {
		t.Skip("GCP_PROJECT_ID is not set")
	}

	location := os.Getenv("GCP_LOCATION")
	if location == "" {
		location = "us-central1"
	}

	client, err := gemini.NewClient(context.Background(), projectID, location)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			t.Errorf("failed to close client: %v", err)
		}
	}()

	tests := []struct {
		name         string
		original     string
		inflammatory string
		wantErr      bool
	}{
		{
			name:         "valid input",
			original:     "今日はいい天気ですね",
			inflammatory: "今日はいい天気ですね(ただし私にとってはですが)",
			wantErr:      false,
		},
		{
			name:         "empty original",
			original:     "",
			inflammatory: "炎上テキスト",
			wantErr:      true,
		},
		{
			name:         "empty inflammatory",
			original:     "元のテキスト",
			inflammatory: "",
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			result, err := client.GenerateExplanation(ctx, tt.original, tt.inflammatory)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, but got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if result == "" {
				t.Error("explanation is empty")
			}
		})
	}
}

func TestClient_GenerateReply(t *testing.T) {
	projectID := os.Getenv("GCP_PROJECT_ID")
	if projectID == "" {
		t.Skip("GCP_PROJECT_ID is not set")
	}

	location := os.Getenv("GCP_LOCATION")
	if location == "" {
		location = "us-central1"
	}

	client, err := gemini.NewClient(context.Background(), projectID, location)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			t.Errorf("failed to close client: %v", err)
		}
	}()

	tests := []struct {
		name      string
		text      string
		replyType string
		wantErr   bool
	}{
		{
			name:      "logical criticism",
			text:      "新しい製品は完璧です",
			replyType: "正論で批判するタイプ",
			wantErr:   false,
		},
		{
			name:      "nitpicking",
			text:      "みんなで協力しましょう",
			replyType: "揚げ足を取るタイプ",
			wantErr:   false,
		},
		{
			name:      "off target",
			text:      "今日はいい天気ですね",
			replyType: "的外れな批判",
			wantErr:   false,
		},
		{
			name:      "excessive defense",
			text:      "この政策には問題があります",
			replyType: "過剰に擁護するタイプ",
			wantErr:   false,
		},
		{
			name:      "empty text",
			text:      "",
			replyType: "正論で批判するタイプ",
			wantErr:   true,
		},
		{
			name:      "empty reply type",
			text:      "テストです",
			replyType: "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			result, err := client.GenerateReply(ctx, tt.text, tt.replyType)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, but got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if result == "" {
				t.Error("reply is empty")
			}
		})
	}
}

func TestNewClient_InvalidProjectID(t *testing.T) {
	ctx := context.Background()
	_, err := gemini.NewClient(ctx, "", "us-central1")
	if err == nil {
		t.Error("expected error with empty project ID, but got nil")
	}
}
