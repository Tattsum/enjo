package image

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/oauth2/google"
)

// ImagenRequest represents the request to Imagen API
type ImagenRequest struct {
	Instances  []ImagenInstance `json:"instances"`
	Parameters ImagenParameters `json:"parameters"`
}

// ImagenInstance represents an instance in the request
type ImagenInstance struct {
	Prompt string `json:"prompt"`
}

// ImagenParameters represents parameters for image generation
type ImagenParameters struct {
	SampleCount     int    `json:"sampleCount"`
	AspectRatio     string `json:"aspectRatio,omitempty"`
	NegativePrompt  string `json:"negativePrompt,omitempty"`
	SampleImageSize string `json:"sampleImageSize,omitempty"`
}

// ImagenResponse represents the response from Imagen API
type ImagenResponse struct {
	Predictions []ImagenPrediction `json:"predictions"`
}

// ImagenPrediction represents a prediction in the response
type ImagenPrediction struct {
	BytesBase64Encoded string `json:"bytesBase64Encoded"`
	MimeType           string `json:"mimeType"`
}

// generateImageViaREST generates an image using Imagen REST API
func (c *Client) generateImageViaREST(ctx context.Context, prompt string, opts *imageOptions) ([]byte, error) {
	// Get OAuth2 token
	creds, err := google.FindDefaultCredentials(ctx, "https://www.googleapis.com/auth/cloud-platform")
	if err != nil {
		return nil, fmt.Errorf("failed to get credentials: %w", err)
	}

	token, err := creds.TokenSource.Token()
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	// Build the API endpoint
	endpoint := fmt.Sprintf(
		"https://%s-aiplatform.googleapis.com/v1/projects/%s/locations/%s/publishers/google/models/%s:predict",
		c.location,
		c.projectID,
		c.location,
		defaultImageModel,
	)

	// Build the request payload
	request := ImagenRequest{
		Instances: []ImagenInstance{
			{Prompt: prompt},
		},
		Parameters: ImagenParameters{
			SampleCount:     defaultSampleCount,
			AspectRatio:     opts.aspectRatio,
			NegativePrompt:  "blurry, low quality, distorted, watermark, text",
			SampleImageSize: "1024",
		},
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check for errors
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	// Parse the response
	var imagenResp ImagenResponse
	if err := json.Unmarshal(body, &imagenResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Extract image data
	if len(imagenResp.Predictions) == 0 {
		return nil, errors.New("no predictions in response")
	}

	// Image is base64 encoded, but we'll return it as-is for now
	// The caller will handle base64 decoding
	return []byte(imagenResp.Predictions[0].BytesBase64Encoded), nil
}
