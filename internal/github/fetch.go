package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// HTTPClient is a real GitHub API client.
type HTTPClient struct {
	BaseURL string
	Token   string
	HTTP    *http.Client
}

// NewHTTPClient creates a GitHub client using GITHUB_TOKEN or gh CLI.
func NewHTTPClient() *HTTPClient {
	token, _ := Token()
	return &HTTPClient{
		BaseURL: "https://api.github.com",
		Token:   token,
		HTTP:    http.DefaultClient,
	}
}

func (c *HTTPClient) request(ctx context.Context, method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}
	return c.HTTP.Do(req)
}

// FetchPR retrieves PR metadata.
func (c *HTTPClient) FetchPR(ctx context.Context, owner, repo string, number int) (PR, error) {
	var pr PR
	resp, err := c.request(ctx, http.MethodGet, fmt.Sprintf("%s/repos/%s/%s/pulls/%d", c.BaseURL, owner, repo, number), nil)
	if err != nil {
		return pr, err
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode != http.StatusOK {
		return pr, fmt.Errorf("github %d", resp.StatusCode)
	}
	if err := json.NewDecoder(resp.Body).Decode(&pr); err != nil {
		return pr, err
	}
	return pr, nil
}

// FetchDiff retrieves the unified diff for the PR.
func (c *HTTPClient) FetchDiff(ctx context.Context, owner, repo string, number int) (string, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/pulls/%d", c.BaseURL, owner, repo, number)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/vnd.github.diff")
	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close() //nolint:errcheck

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("github %d", resp.StatusCode)
	}
	return string(data), nil
}

// FetchFiles retrieves changed files.
func (c *HTTPClient) FetchFiles(ctx context.Context, owner, repo string, number int) ([]PRFile, error) {
	resp, err := c.request(ctx, http.MethodGet, fmt.Sprintf("%s/repos/%s/%s/pulls/%d/files?per_page=100", c.BaseURL, owner, repo, number), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github %d", resp.StatusCode)
	}
	var files []PRFile
	if err := json.NewDecoder(resp.Body).Decode(&files); err != nil {
		return nil, err
	}
	return files, nil
}

// FetchComments retrieves existing PR comments.
func (c *HTTPClient) FetchComments(ctx context.Context, owner, repo string, number int) ([]PRComment, error) {
	resp, err := c.request(ctx, http.MethodGet, fmt.Sprintf("%s/repos/%s/%s/issues/%d/comments?per_page=100", c.BaseURL, owner, repo, number), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github %d", resp.StatusCode)
	}
	var comments []PRComment
	if err := json.NewDecoder(resp.Body).Decode(&comments); err != nil {
		return nil, err
	}
	return comments, nil
}
