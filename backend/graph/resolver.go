package graph

import (
	"context"

	"github.com/Tattsum/enjo/backend/twitter"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// GeminiClient is the interface for Gemini API client
type GeminiClient interface {
	GenerateInflammatoryText(ctx context.Context, original string, level int) (string, error)
	GenerateExplanation(ctx context.Context, original, inflammatory string) (string, error)
	GenerateReply(ctx context.Context, text, replyType string) (string, error)
}

// TwitterClient is the interface for Twitter API client
type TwitterClient interface {
	PostTweet(ctx context.Context, text string, options ...twitter.TweetOption) (*twitter.TweetResult, error)
}

// Resolver is the root resolver for GraphQL
type Resolver struct {
	geminiClient  GeminiClient
	twitterClient TwitterClient
}

// NewResolver creates a new Resolver with dependencies
func NewResolver(geminiClient GeminiClient, twitterClient TwitterClient) *Resolver {
	return &Resolver{
		geminiClient:  geminiClient,
		twitterClient: twitterClient,
	}
}
