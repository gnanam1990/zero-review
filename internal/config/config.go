package config

// Options holds CLI-driven review options.
type Options struct {
	PRURL         string
	Provider      string
	Model         string
	Strict        bool
	SecurityOnly  bool
	TestsOnly     bool
	NoPost        bool
	MaxFindings   int
	ConfidenceMin int
}

// DefaultOptions returns conservative defaults.
func DefaultOptions() Options {
	return Options{
		Provider:      "kimi",
		Model:         "",
		MaxFindings:   50,
		ConfidenceMin: 75,
	}
}
