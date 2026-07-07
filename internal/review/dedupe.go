package review

// Deduplicate removes findings that are too similar, keeping the higher-confidence one.
func Deduplicate(findings []Finding) []Finding {
	var out []Finding
	for _, f := range findings {
		duplicate := false
		for i, existing := range out {
			if similar(f, existing) {
				if f.Confidence > existing.Confidence {
					out[i] = f
				}
				duplicate = true
				break
			}
		}
		if !duplicate {
			out = append(out, f)
		}
	}
	return out
}

func similar(a, b Finding) bool {
	if a.FilePath != b.FilePath {
		return false
	}
	if a.LineStart != b.LineStart {
		return false
	}
	if a.Category != b.Category {
		return false
	}
	if sameFirstWords(a.Title, b.Title, 3) {
		return true
	}
	if sameFirstWords(a.Explanation, b.Explanation, 5) {
		return true
	}
	return false
}

func sameFirstWords(a, b string, n int) bool {
	aw := words(a, n)
	bw := words(b, n)
	if len(aw) < n || len(bw) < n {
		return false
	}
	for i := 0; i < n; i++ {
		if aw[i] != bw[i] {
			return false
		}
	}
	return true
}

func words(s string, max int) []string {
	var out []string
	for _, w := range splitAny(s, " \t\n.,;:!?()[]{}\"'") {
		if w != "" {
			out = append(out, w)
			if len(out) >= max {
				break
			}
		}
	}
	return out
}

func splitAny(s, seps string) []string {
	last := 0
	var out []string
	for i, c := range s {
		if stringsContainsRune(seps, c) {
			out = append(out, s[last:i])
			last = i + 1
		}
	}
	out = append(out, s[last:])
	return out
}

func stringsContainsRune(s string, r rune) bool {
	for _, c := range s {
		if c == r {
			return true
		}
	}
	return false
}
