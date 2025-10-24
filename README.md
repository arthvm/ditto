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
	- **Gemini**: set `GOOGLE_API_KEY` in your environment or `.env` file.
	- **Ollama**: run an Ollama server; configure `OLLAMA_HOST` if it is not `http://localhost:11434`.

## Configuration

Ditto loads environment variables from the shell and from a `.env` file in your project (thanks to `godotenv`). The most relevant knobs are:

| Variable | Required for | Description |
| --- | --- | --- |
| `GOOGLE_API_KEY` | Gemini | API key created in [Google AI Studio](https://aistudio.google.com/app/apikey). |
| `OLLAMA_HOST` | Ollama (optional) | Override the Ollama server URL; defaults to `http://localhost:11434`. |

All commands accept a few shared flags:

- `--provider`: select the LLM provider (default: `gemini`).
- `--prompt`: add extra natural-language context for the model.
- `--issues`: repeatable flag for issue IDs; they show up in commit footers and PR bodies. Example: `--issues 123 --issues PROJ-42`.
- `--base` / `--head`: branch pair used by `ditto pr` (defaults to `main` → current branch).

## Usage

### Generate commits

```sh
# generate a message from staged changes and open it in your editor
ditto commit

# include all tracked changes and amend the previous commit
ditto commit --all --amend

# point Ditto to another provider and add more guidance
ditto commit --provider gemini-flash --prompt "Component: auth; focus on UX copy" --issues 123
```

What happens:

1. Ditto inspects your diff (`git diff --staged` by default).
2. The selected provider composes a Conventional Commit.
3. Ditto runs `git commit -em <message>` so you can tweak it before saving.

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
- Honors `.github/pull_request_template.md`, `docs/pull_request_template.md`, or `PULL_REQUEST_TEMPLATE.md` unless `--no-template` is set.
- Calls `gh pr create` with the generated title/body and opens your editor by default for a final review.

> **Note:** Ollama currently supports commit generation only. Stick with a Gemini provider when running `ditto pr`.

## Providers

| Provider | Alias | Commit support | PR support | Notes |
| --- | --- | --- | --- | --- |
| Gemini 2.5 Pro | `gemini` (default) | ✅ | ✅ | Highest quality results, requires `GOOGLE_API_KEY`. |
| Gemini 2.5 Flash | `gemini-flash` | ✅ | ✅ | Faster/cheaper than Pro, good defaults for most repos. |
| Gemini 2.5 Flash Lite | `gemini-flash-lite` | ✅ | ✅ | Budget-friendly option for smaller diffs. |
| Ollama Git Commit Message | `ollama` | ✅ | ⛔️ | Works with a local Ollama server; ideal for offline commit generation. |

Switch providers on the fly with `--provider`. If a provider lacks a capability (e.g., Ollama PRs), Ditto will surface `not supported with this model`.

## Troubleshooting

- **No staged changes**: run `git add` (or try `--all`) before invoking `ditto commit`.
- **Authentication errors**: ensure `GOOGLE_API_KEY` is exported or stored in `.env`, and you have access to Gemini 2.5 models.
- **`gh` command not found**: install the [GitHub CLI](https://cli.github.com/) and authenticate with `gh auth login`.
- **Ollama connection refused**: start the Ollama server and verify `OLLAMA_HOST`.

## Contributing

We love contributions! To get started:

1. Fork the repository and create a branch (`git checkout -b feature/my-idea`).
2. Make your changes and use `ditto commit` to craft the message.
3. Push and open a pull request. The `ditto pr` command can help you draft it.

Issues and feature requests are also welcome — feel free to open a discussion :)

## License

Distributed under the MIT License. See [LICENSE](LICENSE) for details.
