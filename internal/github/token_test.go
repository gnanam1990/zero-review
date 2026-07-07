package github

import (
	"os"
	"testing"
)

func TestTokenFromEnv(t *testing.T) {
	t.Setenv("GITHUB_TOKEN", "env-token")
	tok, err := Token()
	if err != nil {
		t.Fatalf("Token error: %v", err)
	}
	if tok != "env-token" {
		t.Fatalf("expected env-token, got %q", tok)
	}
}

func TestTokenFallback(t *testing.T) {
	os.Unsetenv("GITHUB_TOKEN")
	tok, err := Token()
	if err != nil {
		t.Fatalf("expected gh fallback to succeed: %v", err)
	}
	if tok == "" {
		t.Fatalf("expected non-empty token from gh fallback")
	}
}
