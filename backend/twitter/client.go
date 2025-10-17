package twitter

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

const (
	// MaxTweetLength is the maximum length of a tweet in characters
	MaxTweetLength = 280
	// MediaUploadURL is the Twitter Media Upload API endpoint
	MediaUploadURL = "https://upload.twitter.com/1.1/media/upload.json"
)

// Client represents a Twitter API client
type Client struct {
	apiKey            string
	apiSecret         string
	accessToken       string
	accessTokenSecret string
	twitterClient     *twitter.Client
	httpClient        *http.Client // OAuth1-authenticated HTTP client
	mockMode          bool         // If true, use mock responses for testing
}

// TweetResult represents the result of posting a tweet
type TweetResult struct {
	ID  string
	URL string
}

// mediaUploadResponse represents the response from Twitter Media Upload API
type mediaUploadResponse struct {
	MediaID       int64  `json:"media_id"`
	MediaIDString string `json:"media_id_string"`
	Size          int    `json:"size"`
	ExpiresAfter  int    `json:"expires_after_secs"`
}

// NewClient creates a new Twitter API client
func NewClient(apiKey, apiSecret, accessToken, accessTokenSecret string) (*Client, error) {
	if apiKey == "" || apiSecret == "" || accessToken == "" || accessTokenSecret == "" {
		return nil, errors.New("all Twitter API credentials are required")
	}

	// Check if we're in test/mock mode (using test credentials)
	mockMode := apiKey == "test-api-key" || apiKey == "test-key"

	var twitterClient *twitter.Client
	var httpClient *http.Client
	if !mockMode {
		// Create OAuth1 config
		config := oauth1.NewConfig(apiKey, apiSecret)
		token := oauth1.NewToken(accessToken, accessTokenSecret)
		httpClient = config.Client(oauth1.NoContext, token)

		// Create Twitter client
		twitterClient = twitter.NewClient(httpClient)
	}

	return &Client{
		apiKey:            apiKey,
		apiSecret:         apiSecret,
		accessToken:       accessToken,
		accessTokenSecret: accessTokenSecret,
		twitterClient:     twitterClient,
		httpClient:        httpClient,
		mockMode:          mockMode,
	}, nil
}

// PostTweet posts a tweet to Twitter
func (c *Client) PostTweet(_ context.Context, text string, options ...TweetOption) (*TweetResult, error) {
	// Validate input
	if text == "" {
		return nil, errors.New("tweet text cannot be empty")
	}

	// Check character limit (considering runes for proper Unicode counting)
	if len([]rune(text)) > MaxTweetLength {
		return nil, errors.New("tweet text exceeds 280 characters")
	}

	// Apply options
	opts := &tweetOptions{}
	for _, opt := range options {
		opt(opts)
	}

	// Build final tweet text
	finalText := c.buildTweetText(text, opts.addHashtag, opts.addDisclaimer)

	// Validate final text length
	if len([]rune(finalText)) > MaxTweetLength {
		return nil, errors.New("tweet text exceeds 280 characters after adding options")
	}

	// If in mock mode, return mock result
	if c.mockMode {
		return &TweetResult{
			ID:  "mock-tweet-id",
			URL: fmt.Sprintf("https://twitter.com/user/status/%s", "mock-tweet-id"),
		}, nil
	}

	// Post tweet using Twitter API
	tweet, _, err := c.twitterClient.Statuses.Update(finalText, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to post tweet: %w", err)
	}

	return &TweetResult{
		ID:  tweet.IDStr,
		URL: fmt.Sprintf("https://twitter.com/%s/status/%s", tweet.User.ScreenName, tweet.IDStr),
	}, nil
}

// buildTweetText builds the final tweet text with optional hashtag and disclaimer
func (*Client) buildTweetText(text string, addHashtag, addDisclaimer bool) string {
	finalText := text

	if addHashtag {
		finalText += " #炎上シミュレーター"
	}

	if addDisclaimer {
		finalText += "\n\n※炎上シミュレーターで生成"
	}

	return finalText
}

// TweetOption is a function type for configuring tweet options
type TweetOption func(*tweetOptions)

type tweetOptions struct {
	addHashtag    bool
	addDisclaimer bool
}

// WithHashtag adds the hashtag to the tweet
func WithHashtag() TweetOption {
	return func(opts *tweetOptions) {
		opts.addHashtag = true
	}
}

// WithDisclaimer adds the disclaimer to the tweet
func WithDisclaimer() TweetOption {
	return func(opts *tweetOptions) {
		opts.addDisclaimer = true
	}
}

// uploadMedia uploads image data to Twitter and returns a media ID
func (c *Client) uploadMedia(ctx context.Context, imageData []byte) (string, error) {
	// Validate input
	if len(imageData) == 0 {
		return "", errors.New("image data cannot be empty")
	}

	// If in mock mode, return mock media ID
	if c.mockMode {
		return "mock-media-id-123456789", nil
	}

	// Base64 encode the image data
	encodedData := base64.StdEncoding.EncodeToString(imageData)

	// Prepare form data
	formData := url.Values{}
	formData.Set("media_data", encodedData)
	formData.Set("media_category", "tweet_image")

	// Create HTTP request with context
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, MediaUploadURL, bytes.NewBufferString(formData.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send request using OAuth1-authenticated HTTP client
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to upload media: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	// Check HTTP status
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("media upload failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse JSON response
	var uploadResp mediaUploadResponse
	if err := json.Unmarshal(body, &uploadResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	return uploadResp.MediaIDString, nil
}

// postTweetWithMediaID posts a tweet with an attached media ID
func (c *Client) postTweetWithMediaID(_ context.Context, text string, mediaID string, options ...TweetOption) (*TweetResult, error) {
	// Validate input
	if text == "" {
		return nil, errors.New("tweet text cannot be empty")
	}
	if mediaID == "" {
		return nil, errors.New("media ID cannot be empty")
	}

	// Apply options
	opts := &tweetOptions{}
	for _, opt := range options {
		opt(opts)
	}

	// Build final tweet text
	finalText := c.buildTweetText(text, opts.addHashtag, opts.addDisclaimer)

	// Validate final text length
	if len([]rune(finalText)) > MaxTweetLength {
		return nil, errors.New("tweet text exceeds 280 characters after adding options")
	}

	// If in mock mode, return mock result
	if c.mockMode {
		return &TweetResult{
			ID:  "mock-tweet-id-with-media",
			URL: fmt.Sprintf("https://twitter.com/user/status/%s", "mock-tweet-id-with-media"),
		}, nil
	}

	// Post tweet with media using Twitter API
	params := &twitter.StatusUpdateParams{
		MediaIds: []int64{mustParseMediaID(mediaID)},
	}
	tweet, _, err := c.twitterClient.Statuses.Update(finalText, params)
	if err != nil {
		return nil, fmt.Errorf("failed to post tweet with media: %w", err)
	}

	return &TweetResult{
		ID:  tweet.IDStr,
		URL: fmt.Sprintf("https://twitter.com/%s/status/%s", tweet.User.ScreenName, tweet.IDStr),
	}, nil
}

// PostTweetWithImage posts a tweet with an image
func (c *Client) PostTweetWithImage(ctx context.Context, text string, imageData []byte, options ...TweetOption) (*TweetResult, error) {
	// Validate input
	if text == "" {
		return nil, errors.New("tweet text cannot be empty")
	}
	if len(imageData) == 0 {
		return nil, errors.New("image data cannot be empty")
	}

	// 1. Upload media
	mediaID, err := c.uploadMedia(ctx, imageData)
	if err != nil {
		return nil, fmt.Errorf("failed to upload media: %w", err)
	}

	// 2. Post tweet with media ID
	result, err := c.postTweetWithMediaID(ctx, text, mediaID, options...)
	if err != nil {
		return nil, fmt.Errorf("failed to post tweet: %w", err)
	}

	return result, nil
}

// mustParseMediaID converts a media ID string to int64
// If parsing fails, it returns 0 (which will cause the API call to fail with a proper error)
func mustParseMediaID(mediaID string) int64 {
	id, err := strconv.ParseInt(mediaID, 10, 64)
	if err != nil {
		return 0
	}
	return id
}
