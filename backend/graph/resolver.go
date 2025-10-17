package graph

import "context"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// GeminiClient is the interface for Gemini API client
type GeminiClient interface {
	GenerateInflammatoryText(ctx context.Context, original string, level int) (string, error)
	GenerateExplanation(ctx context.Context, original, inflammatory string) (string, error)
	GenerateReply(ctx context.Context, text, replyType string) (string, error)
}

// Resolver is the root resolver for GraphQL
type Resolver struct {
	geminiClient GeminiClient
}

// NewResolver creates a new Resolver with dependencies
func NewResolver(geminiClient GeminiClient) *Resolver {
	return &Resolver{
		geminiClient: geminiClient,
	}
}
