package gemini

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

const (
	// Model configuration constants
	defaultTemperature    = 0.9
	defaultTopK           = 40
	defaultTopP           = 0.95
	defaultMaxOutputToken = 1024
)

// Client is a Gemini API client for generating inflammatory text and replies
type Client struct {
	client *genai.Client
	model  *genai.GenerativeModel
}

// NewClient creates a new Gemini API client
func NewClient(ctx context.Context, apiKey string) (*Client, error) {
	if apiKey == "" {
		return nil, errors.New("API key is required")
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	model := client.GenerativeModel("gemini-2.0-flash-exp")

	// Configure the model for consistent output
	model.SetTemperature(defaultTemperature)
	model.SetTopK(defaultTopK)
	model.SetTopP(defaultTopP)
	model.SetMaxOutputTokens(defaultMaxOutputToken)

	return &Client{
		client: client,
		model:  model,
	}, nil
}

// Close closes the Gemini client
func (c *Client) Close() error {
	return c.client.Close()
}

// GenerateInflammatoryText generates inflammatory text from the original text
func (c *Client) GenerateInflammatoryText(ctx context.Context, original string, level int) (string, error) {
	// Validate input
	if original == "" {
		return "", errors.New("original text is required")
	}
	if level < 1 || level > 5 {
		return "", fmt.Errorf("level must be between 1 and 5, got %d", level)
	}

	// Build the prompt
	prompt := buildInflammatoryPrompt(original, level)

	// Generate content
	return c.generateContent(ctx, prompt, "no content generated")
}

// GenerateExplanation generates an explanation of why the text is inflammatory
func (c *Client) GenerateExplanation(ctx context.Context, original, inflammatory string) (string, error) {
	// Validate input
	if original == "" {
		return "", errors.New("original text is required")
	}
	if inflammatory == "" {
		return "", errors.New("inflammatory text is required")
	}

	// Build the prompt
	prompt := buildExplanationPrompt(original, inflammatory)

	// Generate content
	return c.generateContent(ctx, prompt, "no explanation generated")
}

// GenerateReply generates a reply based on the reply type
func (c *Client) GenerateReply(ctx context.Context, text, replyType string) (string, error) {
	// Validate input
	if text == "" {
		return "", errors.New("text is required")
	}
	if replyType == "" {
		return "", errors.New("reply type is required")
	}

	// Build the prompt
	prompt := buildReplyPrompt(text, replyType)

	// Generate content
	return c.generateContent(ctx, prompt, "no reply generated")
}

// generateContent is a helper function to generate content from Gemini API
func (c *Client) generateContent(ctx context.Context, prompt, emptyResultMsg string) (string, error) {
	resp, err := c.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %w", err)
	}

	result := extractTextFromResponse(resp)
	if result == "" {
		return "", errors.New(emptyResultMsg)
	}

	return result, nil
}

// buildInflammatoryPrompt builds a prompt for generating inflammatory text
func buildInflammatoryPrompt(original string, level int) string {
	levelDesc := map[int]string{
		1: "少し配慮に欠ける表現",
		2: "誤解を招きやすい表現",
		3: "明確に批判されそうな表現",
		4: "かなり問題がある表現",
		5: "炎上確実な表現",
	}

	prompt := fmt.Sprintf(`あなたは「炎上シミュレーター」です。以下の投稿を、炎上度レベル%d（1-5）で、
誤解されやすい・批判を受けやすい表現に変換してください。

【元の投稿】
%s

【変換ルール】
- レベル1: 少し配慮に欠ける表現
- レベル2: 誤解を招きやすい表現
- レベル3: 明確に批判されそうな表現
- レベル4: かなり問題がある表現
- レベル5: 炎上確実な表現

【今回のレベル】
レベル%d: %s

変換後の投稿のみを出力してください。説明は不要です。`, level, original, level, levelDesc[level])

	return prompt
}

// buildExplanationPrompt builds a prompt for generating an explanation
func buildExplanationPrompt(original, inflammatory string) string {
	prompt := fmt.Sprintf(`以下の2つの投稿を比較して、なぜ変換後の投稿が炎上しやすいのか、
簡潔に説明してください（2-3文程度）。

【元の投稿】
%s

【変換後の投稿】
%s

変換後の投稿が炎上しやすい理由を、具体的に指摘してください。`, original, inflammatory)

	return prompt
}

// buildReplyPrompt builds a prompt for generating a reply
func buildReplyPrompt(text, replyType string) string {
	typeDesc := map[string]string{
		"正論で批判するタイプ": "正論を振りかざして批判する、理屈っぽいリプライを生成してください。",
		"揚げ足を取るタイプ":  "些細な言葉尻や表現の揚げ足を取る、細かいリプライを生成してください。",
		"的外れな批判":     "投稿の本質とは関係ない、的外れな批判をするリプライを生成してください。",
		"過剰に擁護するタイプ": "投稿を過剰に擁護する、盲目的に賛同するリプライを生成してください。",
	}

	instruction := typeDesc[replyType]
	if instruction == "" {
		instruction = "この投稿に対するリプライを生成してください。"
	}

	prompt := fmt.Sprintf(`以下の投稿に対して、%s

【投稿】
%s

リプライ内容のみを出力してください。説明は不要です。
SNSの投稿のような口調で、簡潔に（2-3文程度）生成してください。`, instruction, text)

	return prompt
}

// extractTextFromResponse extracts text content from Gemini API response
func extractTextFromResponse(resp *genai.GenerateContentResponse) string {
	if resp == nil {
		return ""
	}

	var parts []string
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				if txt, ok := part.(genai.Text); ok {
					parts = append(parts, string(txt))
				}
			}
		}
	}

	return strings.TrimSpace(strings.Join(parts, "\n"))
}
