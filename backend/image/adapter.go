package image

import (
	"context"
)

// Adapter adapts the Client to work with the graph.ImageClient interface
type Adapter struct {
	client *Client
}

// NewAdapter creates a new adapter for the image client
func NewAdapter(client *Client) *Adapter {
	return &Adapter{client: client}
}

// GenerateImage generates an image from a prompt and returns the image data
func (a *Adapter) GenerateImage(ctx context.Context, prompt string) ([]byte, error) {
	result, err := a.client.GenerateImage(ctx, prompt)
	if err != nil {
		return nil, err
	}
	return result.ImageData, nil
}
