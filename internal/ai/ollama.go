package ai

import "context"

// OllamaProvider wraps a local Ollama server using the OpenAI-compatible endpoint.
type OllamaProvider struct {
	inner *OpenAIProvider
}

// NewOllamaProvider creates an Ollama provider.
func NewOllamaProvider(model string) *OllamaProvider {
	if model == "" {
		model = "llama3.1"
	}
	return &OllamaProvider{
		inner: &OpenAIProvider{
			ModelName: model,
			BaseURL:   "http://localhost:11434/v1",
			APIKey:    envOrDefault("OLLAMA_API_KEY", "ollama"),
			Client:    defaultHTTPClient(),
		},
	}
}

// Review sends the prompt to Ollama and parses findings.
func (o *OllamaProvider) Review(ctx context.Context, prompt string) ([]ProviderFinding, error) {
	return o.inner.Review(ctx, prompt)
}

// Model returns the configured model.
func (o *OllamaProvider) Model() string { return o.inner.Model() }

// Name returns the provider name.
func (o *OllamaProvider) Name() string { return "ollama" }
