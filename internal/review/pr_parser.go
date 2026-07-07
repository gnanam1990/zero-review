package review

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

var prURLRE = regexp.MustCompile(`^https?://github\.com/(?P<owner>[\w.-]+)/(?P<repo>[\w.-]+)/pull/(?P<pr>\d+)`)

// ParsePRURL extracts owner, repo, and PR number from a GitHub PR URL.
func ParsePRURL(raw string) (owner, repo string, pr int, err error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", "", 0, fmt.Errorf("empty PR URL")
	}
	u, err := url.Parse(raw)
	if err != nil {
		return "", "", 0, fmt.Errorf("invalid URL: %w", err)
	}
	if u.Host != "github.com" {
		return "", "", 0, fmt.Errorf("only github.com URLs are supported")
	}

	m := prURLRE.FindStringSubmatch(raw)
	if m == nil {
		return "", "", 0, fmt.Errorf("URL does not match github.com/<owner>/<repo>/pull/<number>")
	}

	owner = m[1]
	repo = m[2]
	pr, err = strconv.Atoi(m[3])
	if err != nil {
		return "", "", 0, fmt.Errorf("invalid PR number: %w", err)
	}
	return owner, repo, pr, nil
}
