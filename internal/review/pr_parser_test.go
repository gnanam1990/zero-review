package review

import "testing"

func TestParsePRURL(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		wantOwner string
		wantRepo  string
		wantPR    int
		wantErr   bool
	}{
		{
			name:      "simple PR URL",
			url:       "https://github.com/gnanam1990/zero-review/pull/7",
			wantOwner: "gnanam1990",
			wantRepo:  "zero-review",
			wantPR:    7,
		},
		{
			name:      "URL with trailing slash",
			url:       "https://github.com/gnanam1990/zero-review/pull/42/",
			wantOwner: "gnanam1990",
			wantRepo:  "zero-review",
			wantPR:    42,
		},
		{
			name:    "empty URL",
			url:     "",
			wantErr: true,
		},
		{
			name:    "non-github URL",
			url:     "https://gitlab.com/foo/bar/merge_requests/1",
			wantErr: true,
		},
		{
			name:    "wrong path",
			url:     "https://github.com/gnanam1990/zero-review/issues/7",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			owner, repo, pr, err := ParsePRURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParsePRURL(%q) error = %v, wantErr %v", tt.url, err, tt.wantErr)
			}
			if err != nil {
				return
			}
			if owner != tt.wantOwner || repo != tt.wantRepo || pr != tt.wantPR {
				t.Fatalf("ParsePRURL(%q) = %q, %q, %d; want %q, %q, %d",
					tt.url, owner, repo, pr, tt.wantOwner, tt.wantRepo, tt.wantPR)
			}
		})
	}
}
