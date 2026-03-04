# Ditto

[![Go Version](https://img.shields.io/badge/go-1.25-blue.svg)](https://go.dev/dl/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/arthvm/ditto)](https://goreportcard.com/report/github.com/arthvm/ditto)

Ditto is a friendly command-line companion that turns your staged changes into clean Conventional Commits—and can even draft a polished pull request for you. Under the hood it orchestrates Git, Large Language Models (LLMs), and the GitHub CLI so you can stay in the flow.

## Why Ditto?

- Generate high-quality Conventional Commits (including optional issue footers) from the current diff.
- Use the LLM that works best for you: Google Gemini (multiple models) or a local Ollama model. (Feel free to contribute with more providers!)
- Draft PR titles and bodies that follow your template and include the right context.
- Share extra context through prompts or issue references.
- Configure everything via YAML config files, environment variables, or CLI flags.

## Quickstart

```sh
go install github.com/arthvm/ditto@latest

export GOOGLE_API_KEY="<your-gemini-api-key>" # required for Gemini providers

git add .
ditto commit
```

> Prefer binaries? Grab a pre-built release from the [GitHub releases page](https://github.com/arthvm/ditto/releases).

## Prerequisites

- Go 1.25 or newer (only needed if you install with `go install`).
- Git (Ditto shells out to `git` for commits and diffs).
- [GitHub CLI (`gh`)](https://cli.github.com/) for the `ditto pr` workflow.
- An API key or local model for your chosen provider:
	- **Gemini**: set `GOOGLE_API_KEY` in your environment, `.env` file, or config file.
	- **Ollama**: run an Ollama server; configure the host if it is not `http://localhost:11434`.

## Configuration

Ditto uses a layered configuration system. Settings are resolved in the following order (later sources override earlier ones):

1. **Hardcoded defaults**
2. **User config** (`~/.config/ditto/config.yaml`)
3. **Project config** (`.ditto.yaml` in the repository root)
4. **Environment variables**
5. **CLI flags**

Ditto also loads environment variables from a `.env` file in your project (via `godotenv`).

### Config file

Create `~/.config/ditto/config.yaml` for user-wide settings, or `.ditto.yaml` in a repo root for project-specific overrides. Both use the same format:

```yaml
# Provider selection
provider: gemini            # gemini (default), ollama

# Base branch for PR diffs
base_branch: main

# LLM tuning
llm:
  timeout: "2m"             # request timeout (human-readable duration)
  temperature: 0            # 0 means "use provider default"

# Commit settings
commit:
  prompt: |                 # custom system prompt (replaces the default convention block)
    Use imperative mood.
    Keep the subject line under 50 characters.
  edit: true                # open the editor before committing (default: true)

# PR settings
pr:
  prompt: |                 # custom system prompt (replaces the default convention block)
    Focus on user-facing changes.
  template_path: .github/pull_request_template.md  # custom PR template path
  edit: true                # open the editor before creating the PR (default: true)

# Provider-specific settings (each provider has its own model default)
gemini:
  api_key: ""               # Gemini API key (alternative to GOOGLE_API_KEY env var)
  model: gemini-2.5-flash   # default model for Gemini

ollama:
  host: "http://localhost:11434"
  model: "tavernari/git-commit-message"
```

### Environment variables

| Variable | Description |
| --- | --- |
| `GOOGLE_API_KEY` | Gemini API key. |
| `GEMINI_API_KEY` | Alternative Gemini API key (loaded via config). |
| `OLLAMA_HOST` | Override the Ollama server URL (default: `http://localhost:11434`). |
| `OLLAMA_MODEL` | Override the Ollama model name. |
| `DITTO_PROVIDER` | Override the LLM provider. |
| `DITTO_BASE_BRANCH` | Override the base branch for PR diffs. |
| `DITTO_LLM_TIMEOUT` | Override the LLM timeout (e.g. `"2m"`). |
| `DITTO_LLM_TEMPERATURE` | Override the LLM temperature. |
| `DITTO_COMMIT_EDIT` | Set to `false` or `0` to skip the editor on commit. |
| `DITTO_PR_EDIT` | Set to `false` or `0` to skip the editor on PR creation. |

### CLI flags

All commands accept these global flags:

- `--provider`: select the LLM provider (`gemini`, `ollama`).
- `--model`: override the model for the active provider (e.g. `--provider gemini --model gemini-2.5-pro`).
- `--prompt`: add extra natural-language context for the model.
- `--issues`: repeatable flag for issue IDs; they show up in commit footers and PR bodies. Example: `--issues 123 --issues PROJ-42`.

## Usage

### Generate commits

```sh
# generate a message from staged changes and open it in your editor
ditto commit

# include all tracked changes and amend the previous commit
ditto commit --all --amend

# point Ditto to another provider and add more guidance
ditto commit --provider ollama --prompt "Component: auth; focus on UX copy" --issues 123
```

What happens:

1. Ditto inspects your diff (`git diff --staged` by default).
2. The selected provider composes a Conventional Commit.
3. Ditto runs `git commit -em <message>` so you can tweak it before saving (unless `commit.edit` is set to `false`).

Additional flags:

- `--all`, `-a`: include all tracked changes in the diff.
- `--amend`: regenerate the previous commit message.

### Draft pull requests

```sh
# generate a PR title and body from the diff between branches
ditto pr --base main --head feature/api

# ignore repo templates, create a draft, and add extra context
ditto pr --draft --no-template --prompt "Highlight performance improvements" --issues 456
```

Highlights:

- Uses the commit log and diff stats between `--base` and `--head` to craft a PR narrative.
- Honors `.github/pull_request_template.md`, `docs/pull_request_template.md`, or `PULL_REQUEST_TEMPLATE.md` unless `--no-template` is set. You can also set a custom template path via `pr.template_path` in your config.
- Calls `gh pr create` with the generated title/body and opens your editor by default for a final review (unless `pr.edit` is set to `false`).

### Custom prompts

The `commit.prompt` and `pr.prompt` config options let you define custom system prompts that **replace** the default convention block. This is useful for teams with specific commit or PR conventions:

```yaml
commit:
  prompt: |
    Follow our team's commit convention:
    - Use present tense
    - Prefix with the Jira ticket number
    - Keep subject under 72 characters
```

The `--prompt` CLI flag provides **additional** context on top of the system prompt (custom or default).

## Providers

| Provider | Alias | Commit | PR | Default model | Notes |
| --- | --- | --- | --- | --- | --- |
| Gemini | `gemini` (default) | Yes | Yes | `gemini-2.5-flash` | Requires a Gemini API key. |
| Ollama | `ollama` | Yes | Yes | `tavernari/git-commit-message` | Local Ollama server. |

Each provider has its own default model configured in its config section. Use `--model` to override on a per-invocation basis:

```sh
ditto commit --provider gemini --model gemini-2.5-pro
ditto commit --provider ollama --model codellama
```

## Troubleshooting

- **No staged changes**: run `git add` (or try `--all`) before invoking `ditto commit`.
- **Authentication errors**: ensure your Gemini API key is available via environment, `.env`, or config file.
- **`gh` command not found**: install the [GitHub CLI](https://cli.github.com/) and authenticate with `gh auth login`.
- **Ollama connection refused**: start the Ollama server and verify the host setting.
- **PR template not found**: if you set `pr.template_path` in your config, make sure the file exists—Ditto treats a missing explicit template as an error.

## Contributing

We love contributions! To get started:

1. Fork the repository and create a branch (`git checkout -b feature/my-idea`).
2. Make your changes and use `ditto commit` to craft the message.
3. Push and open a pull request. The `ditto pr` command can help you draft it.

Issues and feature requests are also welcome — feel free to open a discussion :)

## License

Distributed under the MIT License. See [LICENSE](LICENSE) for details.
