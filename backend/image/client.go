package image

import (
	"context"
	"errors"
	"fmt"
	"time"

	"cloud.google.com/go/vertexai/genai"
)

const (
	// Default model for image generation (Imagen 2)
	// Note: As of 2024, Imagen 2 is available as "imagegeneration@006"
	// Imagen 3 may not be available in all regions/projects yet
	defaultImageModel = "imagegeneration@002"
	// Default aspect ratio (1:1 for Twitter)
	defaultAspectRatio = "1:1"
	// Default location for Vertex AI
	defaultLocation = "us-central1"
	// Default number of images to generate
	defaultSampleCount = 1
)

// Client is a Vertex AI client for generating images using Imagen
type Client struct {
	client    *genai.Client
	projectID string
	location  string
}

// Result represents the result of image generation
type Result struct {
	ImageData   []byte    // PNG image data
	ImageURL    string    // GCS URL if saved (optional)
	Prompt      string    // The prompt used for generation
	GeneratedAt time.Time // Timestamp of generation
}

// Option is a functional option for image generation
type Option func(*imageOptions)

// imageOptions holds options for image generation
type imageOptions struct {
	style       string
	aspectRatio string
	width       int
	height      int
}

// NewClient creates a new Imagen API client using Application Default Credentials
func NewClient(ctx context.Context, projectID, location string) (*Client, error) {
	if projectID == "" {
		return nil, errors.New("GCP project ID is required")
	}
	if location == "" {
		location = defaultLocation
	}

	client, err := genai.NewClient(ctx, projectID, location)
	if err != nil {
		return nil, fmt.Errorf("failed to create Vertex AI client: %w", err)
	}

	return &Client{
		client:    client,
		projectID: projectID,
		location:  location,
	}, nil
}

// Close closes the Vertex AI client
func (c *Client) Close() error {
	if c.client != nil {
		return c.client.Close()
	}
	return nil
}

// GenerateImage generates an image based on the prompt
func (c *Client) GenerateImage(ctx context.Context, prompt string, options ...Option) (*Result, error) {
	if prompt == "" {
		return nil, errors.New("prompt is required")
	}

	// Apply options
	opts := &imageOptions{
		aspectRatio: defaultAspectRatio,
	}
	for _, opt := range options {
		opt(opts)
	}

	// Add style to the prompt if specified
	enhancedPrompt := prompt
	if opts.style != "" {
		enhancedPrompt = prompt + getStyleHint(opts.style)
	}

	// Generate the image using REST API
	// The current genai SDK doesn't fully support Imagen API
	imageDataBase64, err := c.generateImageViaREST(ctx, enhancedPrompt, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to generate image: %w", err)
	}

	if len(imageDataBase64) == 0 {
		return nil, errors.New("no image data generated")
	}

	result := &Result{
		ImageData:   imageDataBase64, // Base64 encoded
		Prompt:      prompt,
		GeneratedAt: time.Now(),
	}

	return result, nil
}

// getStyleHint returns the style hint for the given style
func getStyleHint(style string) string {
	switch style {
	case "REALISTIC":
		return ", photorealistic, detailed, high quality"
	case "ILLUSTRATION":
		return ", illustration, artistic, colorful"
	case "MEME":
		return ", meme style, funny, internet culture"
	case "DRAMATIC":
		return ", dramatic lighting, cinematic, intense"
	default:
		return ", " + style
	}
}

// WithStyle sets the image style (e.g., "realistic", "illustration", "meme", "dramatic")
func WithStyle(style string) Option {
	return func(opts *imageOptions) {
		opts.style = style
	}
}

// WithAspectRatio sets the aspect ratio (e.g., "1:1", "16:9", "9:16")
func WithAspectRatio(ratio string) Option {
	return func(opts *imageOptions) {
		opts.aspectRatio = ratio
	}
}

// WithSize sets the image dimensions
func WithSize(width, height int) Option {
	return func(opts *imageOptions) {
		opts.width = width
		opts.height = height
	}
}

// NOTE: extractImageDataFromResponse is no longer used
// We now use REST API which returns base64 encoded data directly
