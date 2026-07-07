package github

import "context"

// FakeClient is a test double that returns canned PR data.
type FakeClient struct {
	PRData     PR
	Diff       string
	Files      []PRFile
	Comments   []PRComment
	LastReview *ReviewSubmission
}

// NewFakeClient creates a fake client with simple canned data.
func NewFakeClient() *FakeClient {
	return &FakeClient{
		PRData: PR{
			Title: "Add auth middleware",
			Body:  "This PR adds authentication middleware.",
			User: struct {
				Login string `json:"login"`
			}{Login: "alice"},
			Head: struct {
				Ref string `json:"ref"`
			}{Ref: "feature/auth"},
			Base: struct {
				Ref string `json:"ref"`
			}{Ref: "main"},
		},
		Diff: `diff --git a/cmd/app/main.go b/cmd/app/main.go
new file mode 100644
index 0000000..1111111
--- /dev/null
+++ b/cmd/app/main.go
@@ -0,0 +1,10 @@
+package main
+
+func main() {
+	client := getClient()
+	client.Do(nil)
+}
+
+func getClient() *Client {
+	return nil
+}
`,
		Files: []PRFile{
			{
				Filename:  "cmd/app/main.go",
				Status:    "added",
				Additions: 10,
				Deletions: 0,
				Patch:     "@@ -0,0 +1,10 @@\n+package main\n+...",
			},
			{
				Filename:  "internal/review/engine.go",
				Status:    "modified",
				Additions: 5,
				Deletions: 2,
				Patch:     "@@ -85,5 +85,8 @@...",
			},
		},
	}
}

// FetchPR returns canned PR metadata.
func (f *FakeClient) FetchPR(ctx context.Context, owner, repo string, number int) (PR, error) {
	return f.PRData, nil
}

// FetchDiff returns a canned diff.
func (f *FakeClient) FetchDiff(ctx context.Context, owner, repo string, number int) (string, error) {
	return f.Diff, nil
}

// FetchFiles returns canned files.
func (f *FakeClient) FetchFiles(ctx context.Context, owner, repo string, number int) ([]PRFile, error) {
	return f.Files, nil
}

// FetchComments returns canned comments.
func (f *FakeClient) FetchComments(ctx context.Context, owner, repo string, number int) ([]PRComment, error) {
	return f.Comments, nil
}

// PostReview records the submission and returns nil.
func (f *FakeClient) PostReview(ctx context.Context, owner, repo string, number int, review ReviewSubmission) error {
	f.LastReview = &review
	return nil
}
