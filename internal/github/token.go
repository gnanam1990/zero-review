package github

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Token returns a GitHub token from the environment or the gh CLI.
func Token() (string, error) {
	if t := os.Getenv("GITHUB_TOKEN"); t != "" {
		return t, nil
	}

	// Try gh CLI as a fallback.
	cmd := exec.Command("gh", "auth", "token")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("GITHUB_TOKEN not set and gh auth token failed: %w", err)
	}
	token := strings.TrimSpace(out.String())
	if token == "" {
		return "", fmt.Errorf("GITHUB_TOKEN not set and gh auth token returned empty")
	}
	return token, nil
}
