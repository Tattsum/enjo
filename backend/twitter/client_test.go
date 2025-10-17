package twitter

import (
	"context"
	"testing"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name              string
		apiKey            string
		apiSecret         string
		accessToken       string
		accessTokenSecret string
		wantErr           bool
	}{
		{
			name:              "valid credentials",
			apiKey:            "test-api-key",
			apiSecret:         "test-api-secret",
			accessToken:       "test-access-token",
			accessTokenSecret: "test-access-token-secret",
			wantErr:           false,
		},
		{
			name:              "empty api key",
			apiKey:            "",
			apiSecret:         "test-api-secret",
			accessToken:       "test-access-token",
			accessTokenSecret: "test-access-token-secret",
			wantErr:           true,
		},
		{
			name:              "empty api secret",
			apiKey:            "test-api-key",
			apiSecret:         "",
			accessToken:       "test-access-token",
			accessTokenSecret: "test-access-token-secret",
			wantErr:           true,
		},
		{
			name:              "empty access token",
			apiKey:            "test-api-key",
			apiSecret:         "test-api-secret",
			accessToken:       "",
			accessTokenSecret: "test-access-token-secret",
			wantErr:           true,
		},
		{
			name:              "empty access token secret",
			apiKey:            "test-api-key",
			apiSecret:         "test-api-secret",
			accessToken:       "test-access-token",
			accessTokenSecret: "",
			wantErr:           true,
		},
		{
			name:              "all empty",
			apiKey:            "",
			apiSecret:         "",
			accessToken:       "",
			accessTokenSecret: "",
			wantErr:           true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.apiKey, tt.apiSecret, tt.accessToken, tt.accessTokenSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && client == nil {
				t.Error("Expected non-nil client for valid credentials")
			}
			if tt.wantErr && client != nil {
				t.Error("Expected nil client for invalid credentials")
			}
		})
	}
}

func TestPostTweet_Validation(t *testing.T) {
	client, err := NewClient("test-key", "test-secret", "test-token", "test-token-secret")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	tests := []struct {
		name    string
		text    string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid text",
			text:    "This is a test tweet",
			wantErr: false,
		},
		{
			name:    "empty text",
			text:    "",
			wantErr: true,
			errMsg:  "tweet text cannot be empty",
		},
		{
			name:    "text exceeds 280 characters",
			text:    "This is a very long tweet that exceeds the maximum length allowed by Twitter API. " + string(make([]byte, 250)),
			wantErr: true,
			errMsg:  "tweet text exceeds 280 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			_, err := client.PostTweet(ctx, tt.text)
			if (err != nil) != tt.wantErr {
				t.Errorf("PostTweet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && err.Error() != tt.errMsg {
				t.Errorf("PostTweet() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestPostTweet_WithOptions(t *testing.T) {
	client, err := NewClient("test-key", "test-secret", "test-token", "test-token-secret")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	tests := []struct {
		name       string
		text       string
		addHashtag bool
		disclaimer bool
		wantText   string
	}{
		{
			name:       "with hashtag",
			text:       "Test tweet",
			addHashtag: true,
			disclaimer: false,
			wantText:   "Test tweet #炎上シミュレーター",
		},
		{
			name:       "with disclaimer",
			text:       "Test tweet",
			addHashtag: false,
			disclaimer: true,
			wantText:   "Test tweet\n\n※炎上シミュレーターで生成",
		},
		{
			name:       "with both",
			text:       "Test tweet",
			addHashtag: true,
			disclaimer: true,
			wantText:   "Test tweet #炎上シミュレーター\n\n※炎上シミュレーターで生成",
		},
		{
			name:       "without options",
			text:       "Test tweet",
			addHashtag: false,
			disclaimer: false,
			wantText:   "Test tweet",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			finalText := client.buildTweetText(tt.text, tt.addHashtag, tt.disclaimer)
			if finalText != tt.wantText {
				t.Errorf("buildTweetText() = %v, want %v", finalText, tt.wantText)
			}
		})
	}
}
