package ai

import (
	"context"
	"fmt"
)

// KimiProvider wraps the Kimi Code OpenAI-compatible API.
type KimiProvider struct {
	inner *OpenAIProvider
	key   string
}

// NewKimiProvider creates a Kimi Code provider.
func NewKimiProvider(model string) *KimiProvider {
	if model == "" {
		model = "kimi-for-coding"
	}
	key := envOrDefault("KIMI_API_KEY", "")
	temp := 1.0
	return &KimiProvider{
		inner: &OpenAIProvider{
			ModelName:   model,
			BaseURL:     "https://api.kimi.com/coding/v1",
			APIKey:      key,
			Client:      defaultHTTPClient(),
			Temperature: &temp,
		},
		key: key,
	}
}

// Review sends the prompt to Kimi and parses findings.
func (k *KimiProvider) Review(ctx context.Context, prompt string) ([]ProviderFinding, error) {
	if k.key == "" {
		return nil, fmt.Errorf("KIMI_API_KEY not set")
	}
	return k.inner.Review(ctx, prompt)
}

// Model returns the configured model.
func (k *KimiProvider) Model() string { return k.inner.Model() }

// Name returns the provider name.
func (k *KimiProvider) Name() string { return "kimi" }
