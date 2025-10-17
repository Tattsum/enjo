package graph

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/Tattsum/enjo/backend/graph/model"
	"github.com/Tattsum/enjo/backend/twitter"
)

// MockGeminiClient is a mock implementation of the Gemini client for testing
type MockGeminiClient struct {
	GenerateInflammatoryTextFunc func(ctx context.Context, original string, level int) (string, error)
	GenerateExplanationFunc      func(ctx context.Context, original, inflammatory string) (string, error)
	GenerateReplyFunc            func(ctx context.Context, text, replyType string) (string, error)
	GenerateContentFunc          func(ctx context.Context, prompt string) (string, error)
}

func (m *MockGeminiClient) GenerateInflammatoryText(ctx context.Context, original string, level int) (string, error) {
	if m.GenerateInflammatoryTextFunc != nil {
		return m.GenerateInflammatoryTextFunc(ctx, original, level)
	}
	return "", errors.New("not implemented")
}

func (m *MockGeminiClient) GenerateExplanation(ctx context.Context, original, inflammatory string) (string, error) {
	if m.GenerateExplanationFunc != nil {
		return m.GenerateExplanationFunc(ctx, original, inflammatory)
	}
	return "", errors.New("not implemented")
}

func (m *MockGeminiClient) GenerateReply(ctx context.Context, text, replyType string) (string, error) {
	if m.GenerateReplyFunc != nil {
		return m.GenerateReplyFunc(ctx, text, replyType)
	}
	return "", errors.New("not implemented")
}

func (m *MockGeminiClient) GenerateContent(ctx context.Context, prompt string) (string, error) {
	if m.GenerateContentFunc != nil {
		return m.GenerateContentFunc(ctx, prompt)
	}
	return "", errors.New("not implemented")
}

// MockImageClient is a mock implementation of the Image client for testing
type MockImageClient struct {
	GenerateImageFunc func(ctx context.Context, prompt string) ([]byte, error)
}

func (m *MockImageClient) GenerateImage(ctx context.Context, prompt string) ([]byte, error) {
	if m.GenerateImageFunc != nil {
		return m.GenerateImageFunc(ctx, prompt)
	}
	return nil, errors.New("not implemented")
}

// MockTwitterClient is a mock implementation of the Twitter client for testing
type MockTwitterClient struct {
	PostTweetFunc          func(ctx context.Context, text string) (*twitter.TweetResult, error)
	PostTweetWithImageFunc func(ctx context.Context, text string, imageData []byte) (*twitter.TweetResult, error)
}

func (m *MockTwitterClient) PostTweet(ctx context.Context, text string, _ ...twitter.TweetOption) (*twitter.TweetResult, error) {
	if m.PostTweetFunc != nil {
		return m.PostTweetFunc(ctx, text)
	}
	return nil, errors.New("not implemented")
}

func (m *MockTwitterClient) PostTweetWithImage(ctx context.Context, text string, imageData []byte, _ ...twitter.TweetOption) (*twitter.TweetResult, error) {
	if m.PostTweetWithImageFunc != nil {
		return m.PostTweetWithImageFunc(ctx, text, imageData)
	}
	return nil, errors.New("not implemented")
}

func TestQueryResolver_Health(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{
			name:    "returns OK",
			want:    "OK",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Resolver{}
			resolver := &queryResolver{r}

			got, err := resolver.Health(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("Health() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Health() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMutationResolver_GenerateInflammatoryText(t *testing.T) {
	tests := []struct {
		name       string
		input      model.GenerateInput
		mockText   string
		mockExpl   string
		mockErr    error
		wantText   string
		wantExpl   string
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "successfully generates inflammatory text",
			input: model.GenerateInput{
				OriginalText: "今日はいい天気ですね",
				Level:        3,
			},
			mockText: "今日は最高の天気なのに働いてるやつwww",
			mockExpl: "労働者を見下すような表現になっており、多くの人の反感を買う可能性があります。",
			mockErr:  nil,
			wantText: "今日は最高の天気なのに働いてるやつwww",
			wantExpl: "労働者を見下すような表現になっており、多くの人の反感を買う可能性があります。",
			wantErr:  false,
		},
		{
			name: "returns error when level is out of range",
			input: model.GenerateInput{
				OriginalText: "テスト投稿",
				Level:        6,
			},
			wantErr:    true,
			wantErrMsg: "level must be between 1 and 5",
		},
		{
			name: "returns error when Gemini API fails",
			input: model.GenerateInput{
				OriginalText: "テスト投稿",
				Level:        3,
			},
			mockErr:    errors.New("API error"),
			wantErr:    true,
			wantErrMsg: "failed to generate inflammatory text",
		},
	}

	for _, tt := range tests {
		// capture range variable
		t.Run(tt.name, func(t *testing.T) {
			mockClient := createMockClientForInflammatory(tt.mockErr, tt.mockText, tt.mockExpl)
			r := &Resolver{geminiClient: mockClient}
			resolver := &mutationResolver{r}

			got, err := resolver.GenerateInflammatoryText(context.Background(), tt.input)

			assertInflammatoryTextResult(t, got, err, tt.wantErr, tt.wantErrMsg, tt.wantText, tt.wantExpl)
		})
	}
}

func createMockClientForInflammatory(mockErr error, mockText, mockExpl string) *MockGeminiClient {
	return &MockGeminiClient{
		GenerateInflammatoryTextFunc: func(_ context.Context, _ string, _ int) (string, error) {
			if mockErr != nil {
				return "", mockErr
			}
			return mockText, nil
		},
		GenerateExplanationFunc: func(_ context.Context, _, _ string) (string, error) {
			if mockErr != nil {
				return "", mockErr
			}
			return mockExpl, nil
		},
	}
}

func assertInflammatoryTextResult(t *testing.T, got *model.GenerateResult, err error, wantErr bool, wantErrMsg, wantText, wantExpl string) {
	t.Helper()

	if (err != nil) != wantErr {
		t.Errorf("GenerateInflammatoryText() error = %v, wantErr %v", err, wantErr)
		return
	}

	if wantErr {
		if err != nil && wantErrMsg != "" && !strings.Contains(err.Error(), wantErrMsg) {
			t.Errorf("GenerateInflammatoryText() error = %v, want error containing %v", err, wantErrMsg)
		}
		return
	}

	if got == nil {
		t.Fatal("GenerateInflammatoryText() returned nil result")
	}

	if got.InflammatoryText != wantText {
		t.Errorf("GenerateInflammatoryText().InflammatoryText = %v, want %v", got.InflammatoryText, wantText)
	}

	if got.Explanation != nil && *got.Explanation != wantExpl {
		t.Errorf("GenerateInflammatoryText().Explanation = %v, want %v", *got.Explanation, wantExpl)
	}
}

func TestMutationResolver_GenerateReplies(t *testing.T) {
	tests := []struct {
		name       string
		text       string
		mockReply  string
		mockErr    error
		wantCount  int
		wantErr    bool
		wantErrMsg string
	}{
		{
			name:      "successfully generates 4 replies",
			text:      "炎上しそうな投稿",
			mockReply: "これはテストリプライです",
			mockErr:   nil,
			wantCount: 4,
			wantErr:   false,
		},
		{
			name:       "returns error when text is empty",
			text:       "",
			wantErr:    true,
			wantErrMsg: "text is required",
		},
		{
			name:       "returns error when Gemini API fails",
			text:       "テスト投稿",
			mockErr:    errors.New("API error"),
			wantErr:    true,
			wantErrMsg: "failed to generate reply",
		},
	}

	for _, tt := range tests {
		// capture range variable
		t.Run(tt.name, func(t *testing.T) {
			mockClient := createMockClientForReplies(tt.mockErr, tt.mockReply)
			r := &Resolver{geminiClient: mockClient}
			resolver := &mutationResolver{r}

			got, err := resolver.GenerateReplies(context.Background(), tt.text)

			assertRepliesResult(t, got, err, tt.wantErr, tt.wantErrMsg, tt.wantCount)
		})
	}
}

func createMockClientForReplies(mockErr error, mockReply string) *MockGeminiClient {
	return &MockGeminiClient{
		GenerateReplyFunc: func(_ context.Context, _, _ string) (string, error) {
			if mockErr != nil {
				return "", mockErr
			}
			return mockReply, nil
		},
	}
}

func assertRepliesResult(t *testing.T, got []*model.Reply, err error, wantErr bool, wantErrMsg string, wantCount int) {
	t.Helper()

	if (err != nil) != wantErr {
		t.Errorf("GenerateReplies() error = %v, wantErr %v", err, wantErr)
		return
	}

	if wantErr {
		if err != nil && wantErrMsg != "" && !strings.Contains(err.Error(), wantErrMsg) {
			t.Errorf("GenerateReplies() error = %v, want error containing %v", err, wantErrMsg)
		}
		return
	}

	if len(got) != wantCount {
		t.Errorf("GenerateReplies() returned %d replies, want %d", len(got), wantCount)
	}

	validateReplyTypes(t, got)
}

func validateReplyTypes(t *testing.T, replies []*model.Reply) {
	t.Helper()

	replyTypes := make(map[model.ReplyType]bool)
	for _, reply := range replies {
		if reply.ID == "" {
			t.Error("GenerateReplies() reply has empty ID")
		}
		if reply.Content == "" {
			t.Error("GenerateReplies() reply has empty Content")
		}
		replyTypes[reply.Type] = true
	}

	expectedTypes := []model.ReplyType{
		model.ReplyTypeLogicalCriticism,
		model.ReplyTypeNitpicking,
		model.ReplyTypeOffTarget,
		model.ReplyTypeExcessiveDefense,
	}

	for _, expectedType := range expectedTypes {
		if !replyTypes[expectedType] {
			t.Errorf("GenerateReplies() missing reply type %v", expectedType)
		}
	}
}

func TestMutationResolver_GenerateImage(t *testing.T) {
	tests := []struct {
		name          string
		input         model.GenerateImageInput
		mockPrompt    string
		mockImageData []byte
		mockPromptErr error
		mockImageErr  error
		wantPrompt    string
		wantErr       bool
		wantErrMsg    string
	}{
		{
			name: "successfully generates image",
			input: model.GenerateImageInput{
				Text: "炎上しそうな投稿",
			},
			mockPrompt:    "A dramatic scene with flames and social media chaos",
			mockImageData: []byte("fake-image-data"),
			wantPrompt:    "A dramatic scene with flames and social media chaos",
			wantErr:       false,
		},
		{
			name: "returns error when text is empty",
			input: model.GenerateImageInput{
				Text: "",
			},
			wantErr:    true,
			wantErrMsg: "text is required",
		},
		{
			name: "returns error when prompt generation fails",
			input: model.GenerateImageInput{
				Text: "テスト投稿",
			},
			mockPromptErr: errors.New("prompt generation error"),
			wantErr:       true,
			wantErrMsg:    "failed to generate image prompt",
		},
		{
			name: "returns error when image generation fails",
			input: model.GenerateImageInput{
				Text: "テスト投稿",
			},
			mockPrompt:   "Test prompt",
			mockImageErr: errors.New("image generation error"),
			wantErr:      true,
			wantErrMsg:   "failed to generate image",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resolver := createImageGenerationResolver(tt.mockPrompt, tt.mockImageData, tt.mockPromptErr, tt.mockImageErr)
			got, err := resolver.GenerateImage(context.Background(), tt.input)
			assertImageGenerationResult(t, got, err, tt.wantErr, tt.wantErrMsg, tt.wantPrompt)
		})
	}
}

func createImageGenerationResolver(mockPrompt string, mockImageData []byte, mockPromptErr, mockImageErr error) *mutationResolver {
	mockGeminiClient := &MockGeminiClient{
		GenerateContentFunc: func(_ context.Context, _ string) (string, error) {
			if mockPromptErr != nil {
				return "", mockPromptErr
			}
			return mockPrompt, nil
		},
	}

	mockImageClient := &MockImageClient{
		GenerateImageFunc: func(_ context.Context, _ string) ([]byte, error) {
			if mockImageErr != nil {
				return nil, mockImageErr
			}
			return mockImageData, nil
		},
	}

	r := &Resolver{
		geminiClient: mockGeminiClient,
		imageClient:  mockImageClient,
	}
	return &mutationResolver{r}
}

func assertImageGenerationResult(t *testing.T, got *model.GenerateImageResult, err error, wantErr bool, wantErrMsg, wantPrompt string) {
	t.Helper()

	if (err != nil) != wantErr {
		t.Errorf("GenerateImage() error = %v, wantErr %v", err, wantErr)
		return
	}

	if wantErr {
		if err != nil && wantErrMsg != "" && !strings.Contains(err.Error(), wantErrMsg) {
			t.Errorf("GenerateImage() error = %v, want error containing %v", err, wantErrMsg)
		}
		return
	}

	if got == nil {
		t.Fatal("GenerateImage() returned nil result")
	}

	if got.Prompt != wantPrompt {
		t.Errorf("GenerateImage().Prompt = %v, want %v", got.Prompt, wantPrompt)
	}

	if got.ImageURL == "" {
		t.Error("GenerateImage().ImageURL is empty")
	}

	if got.GeneratedAt == "" {
		t.Error("GenerateImage().GeneratedAt is empty")
	}
}
