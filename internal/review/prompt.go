package review

import (
	"fmt"
	"strings"

	"github.com/gnanam1990/zero-review/internal/config"
)

// buildPrompt constructs the review prompt from PR data and options.
func buildPrompt(data PRData, opts config.Options) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Review the following pull request.\n\nTitle: %s\nAuthor: %s\nBranch: %s -> %s\nChanged files: %d (+ %d / - %d)\n\n",
		data.Title, data.Author, data.Branch, data.Base, data.ChangedFiles, data.Additions, data.Deletions,
	)

	filters := []string{}
	if opts.SecurityOnly {
		filters = append(filters, "Only report security findings.")
	}
	if opts.TestsOnly {
		filters = append(filters, "Only report test-related findings.")
	}
	if opts.Strict {
		filters = append(filters, "Use a very high bar: report only issues you are highly confident about.")
	} else {
		filters = append(filters, "Report only high-confidence, useful findings. Avoid nitpicks.")
	}
	if len(filters) > 0 {
		b.WriteString(strings.Join(filters, " ") + "\n\n")
	}

	b.WriteString("Return findings as a JSON array. Each finding must have:\n")
	b.WriteString("  id, title, severity (high|medium|low|info), confidence (0-100), category (security|bug|test|performance|maintainability|docs),\n")
	b.WriteString("  file_path, line_start, line_end, diff_side (LEFT|RIGHT), explanation, suggested_comment, suggested_fix, evidence.\n\n")
	b.WriteString("Diff:\n```diff\n")
	b.WriteString(truncateDiff(data.Diff, 12000))
	b.WriteString("\n```\n")
	return b.String()
}

func truncateDiff(diff string, max int) string {
	if len(diff) <= max {
		return diff
	}
	return diff[:max] + "\n... [diff truncated]"
}
