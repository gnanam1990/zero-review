package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// AnthropicProvider calls the Anthropic Messages API.
type AnthropicProvider struct {
	ModelName string
	APIKey    string
	Client    *http.Client
}

// NewAnthropicProvider creates an Anthropic provider.
func NewAnthropicProvider(model string) *AnthropicProvider {
	if model == "" {
		model = "claude-3-5-sonnet-20241022"
	}
	return &AnthropicProvider{
		ModelName: model,
		APIKey:    os.Getenv("ANTHROPIC_API_KEY"),
		Client:    &http.Client{Timeout: 120 * time.Second},
	}
}

// Review sends the prompt and parses JSON findings.
func (a *AnthropicProvider) Review(ctx context.Context, prompt string) ([]ProviderFinding, error) {
	if a.APIKey == "" {
		return nil, fmt.Errorf("ANTHROPIC_API_KEY not set")
	}

	body := map[string]interface{}{
		"model":      a.ModelName,
		"max_tokens": 4096,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
		"system": "You are a senior code reviewer. Return findings as a JSON array only.",
	}
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.anthropic.com/v1/messages", bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("x-api-key", a.APIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := a.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode != http.StatusOK {
		data, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("anthropic %d: %s", resp.StatusCode, string(data))
	}

	var result struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	if len(result.Content) == 0 {
		return nil, fmt.Errorf("no content returned")
	}

	return parseFindings(result.Content[0].Text)
}

// Model returns the configured model.
func (a *AnthropicProvider) Model() string { return a.ModelName }

// Name returns the provider name.
func (a *AnthropicProvider) Name() string { return "anthropic" }
