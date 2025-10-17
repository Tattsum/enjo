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
	template := fmt.Sprintf(`Create a short English image prompt (max 80 words) for this social media post:

"%s"

Style:
- Colorful, eye-catching illustration
- Fun and playful composition
- Social media friendly, vibrant colors
- Modern digital art style

Output ONLY the English prompt, no explanations.`, inflammatoryText)

	return template
}
