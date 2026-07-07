# zero-review

Interactive terminal AI PR reviewer for GitHub.

## Usage

```bash
# Launch the interactive TUI
zero-review tui

# Review a PR with the default Kimi provider
zero-review review https://github.com/owner/repo/pull/123

# Or pass the URL positionally
zero-review https://github.com/owner/repo/pull/123

# Launch the TUI with realistic mock data (no API keys required)
zero-review review https://github.com/owner/repo/pull/123 --mock

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

## TUI screens

| Screen | What it shows |
|--------|---------------|
| **Welcome** | Start screen with quick actions |
| **PR Input** | Enter a PR URL, pick provider/mode, save/post options |
| **Loading Review** | Animated progress of the review pipeline |
| **Dashboard** | PR summary, risk score, severity counts |
| **Findings** | Filterable list of AI findings |
| **Finding Detail** | Full explanation, evidence, suggested fix |
| **Diff** | Inline diff context for the selected finding |
| **Chat** | Ask the reviewer about a finding or edit its comment |
| **Approval** | Final disposition summary and post controls |
| **Report** | Rendered markdown report preview |
| **Settings** | Provider, mode, confidence threshold, theme |

## Key bindings

### Global

| Key | Action |
|-----|--------|
| `q` / `Ctrl+c` | Quit |
| `esc` / `b` | Back / dismiss modal |
| `?` | Toggle help overlay |
| `/` | Open command bar |
| `o` | Go to dashboard/overview |
| `f` | Go to findings |
| `d` | Go to diff |
| `c` | Go to chat |

### Welcome

| Key | Action |
|-----|--------|
| `Enter` / `s` | Start review |
| `l` | Open last report (mock) |
| `,` / `s` | Open settings |

### PR Input / Settings

| Key | Action |
|-----|--------|
| `Tab` | Next field |
| `Shift+Tab` | Previous field |
| `Enter` | Submit form |
| `Esc` | Cancel |

### Findings

| Key | Action |
|-----|--------|
| `↑` / `↓` | Navigate findings |
| `Enter` | Open finding detail |
| `a` | Approve selected finding |
| `r` | Reject selected finding |
| `e` | Edit comment in chat |
| `c` | Chat about finding |
| `d` | View diff |
| `p` | Go to approval / post screen |
| `x` | Clear filters |

### Finding Detail

| Key | Action |
|-----|--------|
| `a` / `r` | Approve / reject finding |
| `e` / `c` | Edit / chat |
| `d` | View diff |
| `esc` / `b` | Back to findings |

### Chat

| Key | Action |
|-----|--------|
| `Enter` | Send message |
| `Shift+Enter` | Newline |
| `Ctrl+l` | Clear chat |
| `Ctrl+r` | Regenerate reply |
| `esc` | Back |

### Approval

| Key | Action |
|-----|--------|
| `p` | Open post confirmation |
| `m` | Cycle posting mode |
| `s` | Save local report |
| `y` / `n` | Confirm / cancel post |
| `esc` | Back |

## Local reports

Markdown and JSON reports are saved under `~/.zero-review/reports/` automatically when a real review finishes, even if nothing is posted to GitHub.

## Install

```bash
go install github.com/gnanam1990/zero-review/cmd/zero-review@latest
```

## Development

```bash
go build ./...
go test -race ./...
go run ./cmd/zero-review tui
go run ./cmd/zero-review review https://github.com/example/repo/pull/1 --mock
```
