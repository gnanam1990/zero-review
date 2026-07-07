package review

import (
	"context"
	"fmt"
	"sort"

	"github.com/gnanam1990/zero-review/internal/ai"
	"github.com/gnanam1990/zero-review/internal/config"
	"github.com/gnanam1990/zero-review/internal/github"
)

// PRData bundles everything needed to review a PR.
type PRData struct {
	Owner        string
	Repo         string
	PRNumber     int
	Title        string
	Body         string
	Author       string
	Branch       string
	Base         string
	ChangedFiles int
	Additions    int
	Deletions    int
	Commits      int
	Diff         string
	Files        []github.PRFile
	Comments     []github.PRComment
}

// Result is the outcome of a review run.
type Result struct {
	PRData    PRData
	Findings  []Finding
	Model     string
	Provider  string
	Timestamp string
}

// Engine orchestrates fetching, analysis, validation, and sorting.
type Engine struct {
	Client   github.Client
	Provider ai.Provider
	Opts     config.Options
}

// Run reviews a PR from URL to validated findings.
func (e *Engine) Run(ctx context.Context, url string) (Result, error) {
	owner, repo, number, err := ParsePRURL(url)
	if err != nil {
		return Result{}, fmt.Errorf("parse pr url: %w", err)
	}

	pr, err := e.Client.FetchPR(ctx, owner, repo, number)
	if err != nil {
		return Result{}, fmt.Errorf("fetch pr: %w", err)
	}

	diff, err := e.Client.FetchDiff(ctx, owner, repo, number)
	if err != nil {
		return Result{}, fmt.Errorf("fetch diff: %w", err)
	}

	files, err := e.Client.FetchFiles(ctx, owner, repo, number)
	if err != nil {
		return Result{}, fmt.Errorf("fetch files: %w", err)
	}

	comments, err := e.Client.FetchComments(ctx, owner, repo, number)
	if err != nil {
		return Result{}, fmt.Errorf("fetch comments: %w", err)
	}

	data := PRData{
		Owner:        owner,
		Repo:         repo,
		PRNumber:     number,
		Title:        pr.Title,
		Body:         pr.Body,
		Author:       pr.User.Login,
		Branch:       pr.Head.Ref,
		Base:         pr.Base.Ref,
		ChangedFiles: len(files),
		Additions:    sumAdditions(files),
		Deletions:    sumDeletions(files),
		Diff:         diff,
		Files:        files,
		Comments:     comments,
	}

	prompt := buildPrompt(data, e.Opts)
	rawFindings, err := e.Provider.Review(ctx, prompt)
	if err != nil {
		return Result{}, fmt.Errorf("ai review: %w", err)
	}

	validated := ValidateAll(FromProviderFindings(rawFindings), data, e.Opts.ConfidenceMin)
	validated = Deduplicate(validated)
	validated = ApplyFilters(validated, e.Opts)

	sort.SliceStable(validated, func(i, j int) bool {
		if validated[i].Confidence != validated[j].Confidence {
			return validated[i].Confidence > validated[j].Confidence
		}
		return validated[i].Severity < validated[j].Severity
	})

	if len(validated) > e.Opts.MaxFindings {
		validated = validated[:e.Opts.MaxFindings]
	}

	return Result{
		PRData:   data,
		Findings: validated,
		Model:    e.Provider.Model(),
		Provider: e.Provider.Name(),
	}, nil
}

func sumAdditions(files []github.PRFile) int {
	total := 0
	for _, f := range files {
		total += f.Additions
	}
	return total
}

func sumDeletions(files []github.PRFile) int {
	total := 0
	for _, f := range files {
		total += f.Deletions
	}
	return total
}
