package github

import "context"

// Client is the GitHub API abstraction used by the review engine.
type Client interface {
	FetchPR(ctx context.Context, owner, repo string, number int) (PR, error)
	FetchDiff(ctx context.Context, owner, repo string, number int) (string, error)
	FetchFiles(ctx context.Context, owner, repo string, number int) ([]PRFile, error)
	FetchComments(ctx context.Context, owner, repo string, number int) ([]PRComment, error)
	PostReview(ctx context.Context, owner, repo string, number int, review ReviewSubmission) error
}

// ReviewSubmission is the payload for the PR review API.
type ReviewSubmission struct {
	Body     string          `json:"body"`
	Event    string          `json:"event"`
	Comments []ReviewComment `json:"comments,omitempty"`
}

// ReviewComment is one inline review comment.
type ReviewComment struct {
	Path     string `json:"path"`
	Line     int    `json:"line"`
	Side     string `json:"side"`
	Body     string `json:"body"`
	CommitID string `json:"commit_id,omitempty"`
}
