package twitter

import (
	"context"
	"testing"
)

func TestUploadMedia(t *testing.T) {
	tests := []struct {
		name      string
		imageData []byte
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "successful media upload",
			imageData: []byte("fake-image-data"),
			wantErr:   false,
		},
		{
			name:      "empty image data",
			imageData: []byte{},
			wantErr:   true,
			errMsg:    "image data cannot be empty",
		},
		{
			name:      "nil image data",
			imageData: nil,
			wantErr:   true,
			errMsg:    "image data cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &Client{
				apiKey:            "test-api-key",
				apiSecret:         "test-api-secret",
				accessToken:       "test-access-token",
				accessTokenSecret: "test-access-token-secret",
				mockMode:          true,
			}

			ctx := context.Background()
			mediaID, err := client.uploadMedia(ctx, tt.imageData)

			if tt.wantErr {
				if err == nil {
					t.Errorf("uploadMedia() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if tt.errMsg != "" && err.Error() != tt.errMsg {
					t.Errorf("uploadMedia() error = %v, want %v", err.Error(), tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("uploadMedia() unexpected error = %v", err)
				return
			}

			if mediaID == "" {
				t.Error("uploadMedia() returned empty mediaID")
			}
		})
	}
}

func TestPostTweetWithMediaID(t *testing.T) {
	tests := []struct {
		name    string
		text    string
		mediaID string
		options []TweetOption
		wantErr bool
		errMsg  string
	}{
		{
			name:    "successful tweet with media",
			text:    "Test tweet with image",
			mediaID: "123456789",
			wantErr: false,
		},
		{
			name:    "empty text",
			text:    "",
			mediaID: "123456789",
			wantErr: true,
			errMsg:  "tweet text cannot be empty",
		},
		{
			name:    "empty media ID",
			text:    "Test tweet",
			mediaID: "",
			wantErr: true,
			errMsg:  "media ID cannot be empty",
		},
		{
			name:    "with hashtag option",
			text:    "Test tweet",
			mediaID: "123456789",
			options: []TweetOption{WithHashtag()},
			wantErr: false,
		},
		{
			name:    "with disclaimer option",
			text:    "Test tweet",
			mediaID: "123456789",
			options: []TweetOption{WithDisclaimer()},
			wantErr: false,
		},
	}

	//nolint:dupl // Similar structure is intentional for testing different methods
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &Client{
				apiKey:            "test-api-key",
				apiSecret:         "test-api-secret",
				accessToken:       "test-access-token",
				accessTokenSecret: "test-access-token-secret",
				mockMode:          true,
			}

			ctx := context.Background()
			result, err := client.postTweetWithMediaID(ctx, tt.text, tt.mediaID, tt.options...)

			if tt.wantErr {
				if err == nil {
					t.Errorf("postTweetWithMediaID() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if tt.errMsg != "" && err.Error() != tt.errMsg {
					t.Errorf("postTweetWithMediaID() error = %v, want %v", err.Error(), tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("postTweetWithMediaID() unexpected error = %v", err)
				return
			}

			if result == nil {
				t.Error("postTweetWithMediaID() returned nil result")
				return
			}

			if result.ID == "" {
				t.Error("postTweetWithMediaID() returned empty ID")
			}

			if result.URL == "" {
				t.Error("postTweetWithMediaID() returned empty URL")
			}
		})
	}
}

func TestPostTweetWithImage(t *testing.T) {
	tests := []struct {
		name      string
		text      string
		imageData []byte
		options   []TweetOption
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "successful tweet with image",
			text:      "Test tweet with image",
			imageData: []byte("fake-image-data"),
			wantErr:   false,
		},
		{
			name:      "empty text",
			text:      "",
			imageData: []byte("fake-image-data"),
			wantErr:   true,
			errMsg:    "tweet text cannot be empty",
		},
		{
			name:      "empty image data",
			text:      "Test tweet",
			imageData: []byte{},
			wantErr:   true,
			errMsg:    "image data cannot be empty",
		},
		{
			name:      "with both options",
			text:      "Test tweet",
			imageData: []byte("fake-image-data"),
			options:   []TweetOption{WithHashtag(), WithDisclaimer()},
			wantErr:   false,
		},
	}

	//nolint:dupl // Similar structure is intentional for testing different methods
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &Client{
				apiKey:            "test-api-key",
				apiSecret:         "test-api-secret",
				accessToken:       "test-access-token",
				accessTokenSecret: "test-access-token-secret",
				mockMode:          true,
			}

			ctx := context.Background()
			result, err := client.PostTweetWithImage(ctx, tt.text, tt.imageData, tt.options...)

			if tt.wantErr {
				if err == nil {
					t.Errorf("PostTweetWithImage() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if tt.errMsg != "" && err.Error() != tt.errMsg {
					t.Errorf("PostTweetWithImage() error = %v, want %v", err.Error(), tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("PostTweetWithImage() unexpected error = %v", err)
				return
			}

			if result == nil {
				t.Error("PostTweetWithImage() returned nil result")
				return
			}

			if result.ID == "" {
				t.Error("PostTweetWithImage() returned empty ID")
			}

			if result.URL == "" {
				t.Error("PostTweetWithImage() returned empty URL")
			}
		})
	}
}
