package image

import (
	"context"
	"strings"
	"testing"
)

// mockGeminiClient is a mock implementation of gemini client for testing
type mockGeminiClient struct {
	generateContentFunc func(ctx context.Context, prompt string) (string, error)
}

func (m *mockGeminiClient) GenerateContent(ctx context.Context, prompt string) (string, error) {
	if m.generateContentFunc != nil {
		return m.generateContentFunc(ctx, prompt)
	}
	return "A dramatic scene with flames and social media chaos", nil
}

func TestGenerateImagePrompt(t *testing.T) {
	ctx := context.Background()

	t.Run("successful prompt generation", func(t *testing.T) {
		mockClient := &mockGeminiClient{
			generateContentFunc: func(_ context.Context, prompt string) (string, error) {
				// Verify that the prompt contains the inflammatory text
				if !strings.Contains(prompt, "テスト投稿") {
					t.Error("expected prompt to contain inflammatory text")
				}
				return "A dramatic image of fire and controversy on social media", nil
			},
		}

		text := "テスト投稿"
		result, err := GenerateImagePrompt(ctx, mockClient, text)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if result == "" {
			t.Error("expected non-empty result")
		}
	})

	t.Run("error when inflammatory text is empty", func(t *testing.T) {
		mockClient := &mockGeminiClient{}
		_, err := GenerateImagePrompt(ctx, mockClient, "")
		if err == nil {
			t.Fatal("expected error when inflammatory text is empty")
		}
	})

	t.Run("error from gemini client", func(t *testing.T) {
		mockClient := &mockGeminiClient{
			generateContentFunc: func(_ context.Context, _ string) (string, error) {
				return "", context.DeadlineExceeded
			},
		}

		text := "テスト投稿"
		_, err := GenerateImagePrompt(ctx, mockClient, text)
		if err == nil {
			t.Fatal("expected error from gemini client")
		}
	})
}

func TestBuildImagePromptTemplate(t *testing.T) {
	t.Run("builds prompt with inflammatory text", func(t *testing.T) {
		text := "これは炎上しやすい投稿です"
		prompt := buildImagePromptTemplate(text)

		if !strings.Contains(prompt, text) {
			t.Errorf("expected prompt to contain text %q", text)
		}
		if !strings.Contains(prompt, "炎") {
			t.Error("expected prompt to mention fire/flames (炎)")
		}
		if !strings.Contains(prompt, "SNS") {
			t.Error("expected prompt to mention SNS")
		}
	})
}
