package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gnanam1990/zero-review/internal/ai"
	"github.com/gnanam1990/zero-review/internal/config"
	"github.com/gnanam1990/zero-review/internal/github"
	"github.com/gnanam1990/zero-review/internal/report"
	"github.com/gnanam1990/zero-review/internal/review"
	"github.com/gnanam1990/zero-review/internal/tui"
)

func main() {
	opts := config.DefaultOptions()

	// Pre-scan for flags that users naturally place after positional args.
	// Go's flag package stops parsing once it sees the first positional, so
	// we extract these manually first and then let flag.Parse() handle the rest.
	mockFlag := scanBoolFlag("--mock")
	fakeFlag := scanBoolFlag("--fake")

	flag.StringVar(&opts.PRURL, "pr", "", "GitHub PR URL to review")
	flag.StringVar(&opts.PRURL, "url", "", "GitHub PR URL to review (alias)")
	flag.StringVar(&opts.Provider, "provider", opts.Provider, "AI provider (kimi|ollama|anthropic|openai|fake)")
	flag.StringVar(&opts.Model, "model", opts.Model, "AI model override")
	flag.BoolVar(&opts.Strict, "strict", opts.Strict, "Use very high confidence threshold")
	flag.BoolVar(&opts.SecurityOnly, "security-only", opts.SecurityOnly, "Only report security findings")
	flag.BoolVar(&opts.TestsOnly, "tests-only", opts.TestsOnly, "Only report test-related findings")
	flag.BoolVar(&opts.NoPost, "no-post", opts.NoPost, "Never post to GitHub")
	flag.IntVar(&opts.MaxFindings, "max-findings", opts.MaxFindings, "Maximum findings to display")
	flag.IntVar(&opts.ConfidenceMin, "confidence", opts.ConfidenceMin, "Minimum confidence 0-100")
	fake := flag.Bool("fake", false, "Use fake provider + fake GitHub client for demo")
	mock := flag.Bool("mock", false, "Launch TUI with mock data instead of running a real review")
	reportDir := flag.String("report-dir", filepath.Join(os.Getenv("HOME"), ".zero-review", "reports"), "Directory for local reports")
	flag.Parse()

	args := flag.Args()
	*mock = *mock || mockFlag
	*fake = *fake || fakeFlag

	switch {
	case len(args) == 0 || (len(args) == 1 && args[0] == "tui"):
		if err := runTUI(nil); err != nil {
			fmt.Fprintf(os.Stderr, "TUI error: %v\n", err)
			os.Exit(1)
		}
		return

	case len(args) >= 2 && (args[0] == "review" || args[0] == "r"):
		opts.PRURL = args[1]
		if *mock {
			if err := runTUI(tui.MockSession()); err != nil {
				fmt.Fprintf(os.Stderr, "TUI error: %v\n", err)
				os.Exit(1)
			}
			return
		}
		if err := runReview(opts, *fake, *reportDir); err != nil {
			fmt.Fprintf(os.Stderr, "Review failed: %v\n", err)
			os.Exit(1)
		}
		return

	case len(args) == 1 && args[0] == "review":
		fmt.Fprintln(os.Stderr, "Usage: zero-review review <github-pr-url>")
		os.Exit(1)

	default:
		if opts.PRURL == "" {
			fmt.Fprintln(os.Stderr, "Usage: zero-review tui")
			fmt.Fprintln(os.Stderr, "       zero-review review <github-pr-url> [--mock]")
			os.Exit(1)
		}
		// Treat bare URL as review command
		if err := runReview(opts, *fake, *reportDir); err != nil {
			fmt.Fprintf(os.Stderr, "Review failed: %v\n", err)
			os.Exit(1)
		}
	}
}

func runTUI(session *review.ReviewSession) error {
	var p *tea.Program
	if session != nil {
		p = tui.NewAppWithSession(session)
	} else {
		p = tui.NewApp()
	}
	_, err := p.Run()
	return err
}

func runReview(opts config.Options, fake bool, reportDir string) error {
	var gh github.Client
	var provider ai.Provider

	if fake {
		gh = github.NewFakeClient()
		provider = ai.NewFakeProvider()
	} else {
		gh = github.NewHTTPClient()
		if _, err := github.Token(); err != nil {
			return fmt.Errorf("GitHub token not available: %w", err)
		}
		switch opts.Provider {
		case "openai":
			provider = ai.NewOpenAIProvider(opts.Model, "")
		case "kimi":
			provider = ai.NewKimiProvider(opts.Model)
		case "ollama":
			provider = ai.NewOllamaProvider(opts.Model)
		case "anthropic", "":
			opts.Provider = "anthropic"
			provider = ai.NewAnthropicProvider(opts.Model)
		default:
			return fmt.Errorf("unknown provider %q", opts.Provider)
		}
	}

	if opts.Strict {
		opts.ConfidenceMin = 90
	}

	engine := review.Engine{Client: gh, Provider: provider, Opts: opts}
	result, err := engine.Run(context.Background(), opts.PRURL)
	if err != nil {
		return err
	}

	session := result.ToSession(opts)
	if err := runTUI(&session); err != nil {
		return err
	}

	return saveLocalReports(result, reportDir)
}

func saveLocalReports(result review.Result, dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	mdPath, err := report.SaveMarkdown(result, dir)
	if err != nil {
		return err
	}
	jsonPath, err := report.SaveJSON(result, dir)
	if err != nil {
		return err
	}
	fmt.Printf("Report saved:\n  %s\n  %s\n", mdPath, jsonPath)
	return nil
}

// scanBoolFlag searches os.Args for a flag so it works even when placed after
// positional arguments, which Go's flag package does not support by default.
func scanBoolFlag(name string) bool {
	for _, a := range os.Args[1:] {
		if a == name {
			return true
		}
	}
	return false
}
