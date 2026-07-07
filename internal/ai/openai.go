package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func envOrDefault(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

func defaultHTTPClient() *http.Client {
	return &http.Client{Timeout: 120 * time.Second}
}

// OpenAIProvider calls an OpenAI-compatible chat completions endpoint.
type OpenAIProvider struct {
	ModelName   string
	BaseURL     string
	APIKey      string
	Client      *http.Client
	Temperature *float64
}

// NewOpenAIProvider creates an OpenAI provider from env/env/model options.
func NewOpenAIProvider(model, baseURL string) *OpenAIProvider {
	if model == "" {
		model = "gpt-4o"
	}
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}
	return &OpenAIProvider{
		ModelName: model,
		BaseURL:   baseURL,
		APIKey:    os.Getenv("OPENAI_API_KEY"),
		Client:    &http.Client{Timeout: 120 * time.Second},
	}
}

// Review sends the prompt and parses the returned JSON.
func (o *OpenAIProvider) Review(ctx context.Context, prompt string) ([]ProviderFinding, error) {
	if o.APIKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY not set")
	}

	body := map[string]interface{}{
		"model": o.ModelName,
		"messages": []map[string]string{
			{"role": "system", "content": "You are a senior code reviewer. Return findings as a JSON array only."},
			{"role": "user", "content": prompt},
		},
		"temperature": o.temperature(),
	}
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, o.BaseURL+"/chat/completions", bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+o.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := o.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode != http.StatusOK {
		data, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("openai %d: %s", resp.StatusCode, string(data))
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	if len(result.Choices) == 0 {
		return nil, fmt.Errorf("no choices returned")
	}

	return parseFindings(result.Choices[0].Message.Content)
}

// Model returns the configured model.
func (o *OpenAIProvider) Model() string { return o.ModelName }

// Name returns the provider name.
func (o *OpenAIProvider) Name() string { return "openai" }

func (o *OpenAIProvider) temperature() float64 {
	if o.Temperature != nil {
		return *o.Temperature
	}
	return 0.2
}

func parseFindings(content string) ([]ProviderFinding, error) {
	content = strings.TrimSpace(content)
	if strings.HasPrefix(content, "```") {
		content = stripCodeFence(content)
	}
	var findings []ProviderFinding
	if err := json.Unmarshal([]byte(content), &findings); err != nil {
		return nil, fmt.Errorf("parse findings: %w", err)
	}
	return findings, nil
}

func stripCodeFence(s string) string {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "```json")
	s = strings.TrimPrefix(s, "```")
	s = strings.TrimSuffix(s, "```")
	return strings.TrimSpace(s)
}
