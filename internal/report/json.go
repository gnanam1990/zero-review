package report

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gnanam1990/zero-review/internal/review"
)

// SaveJSON writes the full review result as JSON.
func SaveJSON(result review.Result, dir string) (string, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil { //nolint:gosec
		return "", err
	}
	ts := time.Now().UTC().Format("20060102-150405")
	path := filepath.Join(dir, fmt.Sprintf("pr-%d-%s.json", result.PRData.PRNumber, ts))

	out, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", err
	}
	if err := os.WriteFile(path, out, 0o644); err != nil { //nolint:gosec
		return "", err
	}
	return path, nil
}
