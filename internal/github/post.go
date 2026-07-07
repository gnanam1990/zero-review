package github

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// PostReview submits a PR review to GitHub.
func (c *HTTPClient) PostReview(ctx context.Context, owner, repo string, number int, review ReviewSubmission) error {
	payload, err := json.Marshal(review)
	if err != nil {
		return err
	}
	resp, err := c.request(ctx, http.MethodPost, fmt.Sprintf("%s/repos/%s/%s/pulls/%d/reviews", c.BaseURL, owner, repo, number), bytes.NewReader(payload))
	if err != nil {
		return err
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("github post review %d", resp.StatusCode)
	}
	return nil
}
