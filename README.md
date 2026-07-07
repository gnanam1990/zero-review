# zero-review

Interactive terminal AI PR reviewer for GitHub.

## Usage

```bash
# Review a PR with the default Anthropic provider
zero-review review https://github.com/owner/repo/pull/123

# Or pass the URL positionally
zero-review https://github.com/owner/repo/pull/123

# Use the fake provider + fake GitHub client for an offline demo
zero-review review https://github.com/owner/repo/pull/123 --fake

# Only preview findings; never post to GitHub
zero-review review https://github.com/owner/repo/pull/123 --no-post

# Use OpenAI
zero-review review https://github.com/owner/repo/pull/123 --provider openai --model gpt-4o
```

## Required environment

- `GITHUB_TOKEN` for fetching diffs and posting reviews. If unset, zero-review falls back to `gh auth token`.
- `KIMI_API_KEY` when using the default Kimi provider (`--provider kimi`).
- `OLLAMA_API_KEY` when using `--provider ollama` (default: `ollama`).
- `ANTHROPIC_API_KEY` when using `--provider anthropic`.
- `OPENAI_API_KEY` when using `--provider openai`.

## Key bindings

| Screen | Key | Action |
|--------|-----|--------|
| Summary | `Enter` / `s` | Start review |
| Summary | `d` | View diff |
| Findings | `Enter` | Open finding detail |
| Findings | `a` / `r` | Approve / reject selected finding |
| Findings | `e` | Edit comment in chat |
| Findings | `c` | Chat about finding |
| Findings | `p` | Go to approval / post screen |
| Findings | `s` | Save local reports |
| Detail | `a` / `r` | Approve / reject finding |
| Detail | `e` / `c` | Edit / chat |
| Detail | `d` | View diff |
| Approval | `a` / `r` | Approve / reject all |
| Approval | `c` | Clear all dispositions |
| Approval | `Enter` | Post approved findings |
| Anywhere | `q` / `Ctrl+c` | Quit |

## Local reports

Markdown and JSON reports are saved under `~/.zero-review/reports/` automatically when the TUI exits, even if nothing is posted to GitHub.

## Install

```bash
go install github.com/gnanam1990/zero-review/cmd/zero-review@latest
```
