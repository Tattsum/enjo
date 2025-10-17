package graph

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"
)

// stringPtr returns a pointer to a string
func stringPtr(s string) *string {
	return &s
}

// generateImagePromptFromText generates an image prompt from the given text using Gemini
func generateImagePromptFromText(ctx context.Context, client GeminiClient, text string) (string, error) {
	promptTemplate := fmt.Sprintf(`以下の炎上投稿に合わせた、視覚的にインパクトのある画像のプロンプトを生成してください。

【投稿】
%s

【要件】
- 投稿の雰囲気を視覚的に表現
- 炎のモチーフを含める
- SNS映えする構図
- ミーム的な要素
- 日本のネット文化に馴染む表現

画像生成プロンプト（英語）のみを出力してください。説明は不要です。`, text)

	return client.GenerateContent(ctx, promptTemplate)
}

// createImageDataURL creates a data URL from image bytes
func createImageDataURL(imageData []byte) string {
	// Encode image data as base64
	encoded := base64.StdEncoding.EncodeToString(imageData)
	// Return as data URL (assuming PNG format)
	return fmt.Sprintf("data:image/png;base64,%s", encoded)
}

// getCurrentTimestamp returns the current timestamp in RFC3339 format
func getCurrentTimestamp() string {
	return time.Now().Format(time.RFC3339)
}

// extractImageDataFromURL extracts image data from a data URL
func extractImageDataFromURL(dataURL string) ([]byte, error) {
	// Check if it's a data URL
	if !strings.HasPrefix(dataURL, "data:") {
		return nil, errors.New("invalid data URL format")
	}

	// Split at the comma to separate metadata from data
	const numParts = 2
	parts := strings.SplitN(dataURL, ",", numParts)
	if len(parts) != numParts {
		return nil, errors.New("invalid data URL format: missing comma separator")
	}

	// Check if it's base64 encoded
	if !strings.Contains(parts[0], "base64") {
		return nil, errors.New("data URL must be base64 encoded")
	}

	// Decode the base64 data
	imageData, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 data: %w", err)
	}

	return imageData, nil
}
