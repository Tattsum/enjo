package image

import (
	"context"
	"errors"
	"fmt"
)

// GeminiClient is an interface for generating content using Gemini
type GeminiClient interface {
	GenerateContent(ctx context.Context, prompt string) (string, error)
}

// GenerateImagePrompt generates an image generation prompt from inflammatory text using Gemini
func GenerateImagePrompt(ctx context.Context, geminiClient GeminiClient, inflammatoryText string) (string, error) {
	if inflammatoryText == "" {
		return "", errors.New("inflammatory text is required")
	}

	// Build the prompt template
	promptTemplate := buildImagePromptTemplate(inflammatoryText)

	// Use Gemini to generate the image prompt
	imagePrompt, err := geminiClient.GenerateContent(ctx, promptTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to generate image prompt: %w", err)
	}

	return imagePrompt, nil
}

// buildImagePromptTemplate builds a template for generating image prompts
func buildImagePromptTemplate(inflammatoryText string) string {
	template := fmt.Sprintf(`以下の炎上投稿に合わせた、視覚的にインパクトのある画像のプロンプトを生成してください。

【投稿】
%s

【要件】
- 投稿の雰囲気を視覚的に表現
- 炎のモチーフを含める
- SNS映えする構図
- ミーム的な要素
- 日本のネット文化に馴染む表現

画像生成プロンプト（英語）のみを出力してください。説明は不要です。`, inflammatoryText)

	return template
}
