package review

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/gnanam1990/zero-review/internal/github"
)

// ValidateAll filters and validates findings against the PR diff.
func ValidateAll(findings []Finding, data PRData, confidenceMin int) []Finding {
	fileSet := make(map[string]github.PRFile)
	for _, f := range data.Files {
		fileSet[f.Filename] = f
	}

	var out []Finding
	for _, f := range findings {
		if f.ID == "" {
			f.ID = fmt.Sprintf("F%03d", len(out)+1)
		}
		if !isValidSeverity(f.Severity) || !isValidCategory(f.Category) {
			continue
		}
		if f.Confidence < confidenceMin {
			continue
		}
		pf, ok := fileSet[f.FilePath]
		if !ok {
			// Try matching by basename in case AI returns relative paths.
			pf = matchByBasename(fileSet, f.FilePath)
			if pf.Filename == "" {
				continue
			}
			f.FilePath = pf.Filename
		}
		if !lineInFile(f, pf, data.Diff) {
			continue
		}
		if f.LineEnd < f.LineStart {
			f.LineEnd = f.LineStart
		}
		if f.Status == "" {
			f.Status = FindingStatusPending
		}
		if f.DiffSide == "" {
			f.DiffSide = "RIGHT"
		}
		out = append(out, f)
	}
	return out
}

func isValidSeverity(s Severity) bool {
	switch s {
	case SeverityHigh, SeverityMedium, SeverityLow, SeverityInfo:
		return true
	}
	return false
}

func isValidCategory(c Category) bool {
	switch c {
	case CategorySecurity, CategoryBug, CategoryTest, CategoryPerformance, CategoryMaintainability, CategoryDocs:
		return true
	}
	return false
}

func matchByBasename(fileSet map[string]github.PRFile, path string) github.PRFile {
	base := filepath.Base(path)
	for _, f := range fileSet {
		if filepath.Base(f.Filename) == base {
			return f
		}
	}
	return github.PRFile{}
}

func lineInFile(f Finding, pf github.PRFile, diff string) bool {
	if f.LineStart <= 0 {
		return false
	}
	patch := pf.Patch
	if patch == "" {
		patch = extractFileDiff(diff, pf.Filename)
	}
	return lineInHunk(f, patch)
}

func extractFileDiff(diff, filename string) string {
	lines := strings.Split(diff, "\n")
	var out []string
	inFile := false
	for _, line := range lines {
		if strings.HasPrefix(line, "diff --git") {
			inFile = strings.Contains(line, filename)
		}
		if inFile {
			out = append(out, line)
		}
	}
	return strings.Join(out, "\n")
}

func lineInHunk(f Finding, patch string) bool {
	// Capture: oldStart, oldCount, newStart, newCount. Counts default to 1.
	re := regexp.MustCompile(`@@ -(\d+)(?:,(\d+))? \+(\d+)(?:,(\d+))? @@`)
	for _, m := range re.FindAllStringSubmatch(patch, -1) {
		oldStart := atoi(m[1])
		oldCount := 1
		if m[2] != "" {
			oldCount = atoi(m[2])
		}
		newStart := atoi(m[3])
		newCount := 1
		if m[4] != "" {
			newCount = atoi(m[4])
		}

		side := f.DiffSide
		if side == "" {
			side = "RIGHT"
		}
		if side == "RIGHT" && f.LineStart >= newStart && f.LineStart < newStart+newCount {
			return true
		}
		if side == "LEFT" && f.LineStart >= oldStart && f.LineStart < oldStart+oldCount {
			return true
		}
	}
	return false
}

func atoi(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}
